const nameLength = 15;
const textEncoder = new TextEncoder();
const textDecoder = new TextDecoder();

export const action = {
  start: Uint8Array.of(0x04),
  rejoin: Uint8Array.of(0x06),
  cooperative: Uint8Array.of(0x09, 1, ...textEncoder.encode("easy")),
  reveal: (x, y) => Uint8Array.of(0x20, x, y),
  flag: (x, y) => Uint8Array.of(0x21, x, y),
  cursor: (id, x, y) => Uint8Array.of(0x22, id, x, y),
};

export function decodeFrame(buffer) {
  const bytes = new Uint8Array(buffer);
  const view = new DataView(buffer);
  const read16 = (offset) => view.getUint16(offset);
  const requireBytes = (length) => {
    if (bytes.length < length) throw new Error("服务端消息长度不足");
  };
  requireBytes(1);
  switch (bytes[0]) {
    case 0x14:
      requireBytes(2);
      return { type: "welcome", id: bytes[1] };
    case 0x15: {
      requireBytes(3);
      const count = bytes[2];
      const end = 3 + count * (1 + nameLength);
      requireBytes(end + 3);
      const players = [];
      for (let start = 3; start < end; start += 1 + nameLength) {
        const name = textDecoder.decode(bytes.subarray(start + 1, start + 1 + nameLength)).replace(/\x00+$/, "");
        players.push({ id: bytes[start], name });
      }
      return { type: "lobby", host: bytes[1], players, mode: bytes[end + 2] };
    }
    case 0x30: {
      requireBytes(6);
      const width = bytes[1];
      const height = bytes[2];
      const scoreCount = bytes[5];
      const start = 6 + scoreCount * 2;
      requireBytes(start + width * height);
      return {
        type: "board",
        width,
        height,
        flagsLeft: read16(3),
        scores: Array.from({ length: scoreCount }, (_, index) => read16(6 + index * 2)),
        cells: bytes.subarray(start, start + width * height),
      };
    }
    case 0x34: {
      requireBytes(2);
      const count = bytes[1];
      requireBytes(2 + count * 3);
      const cursors = [];
      for (let start = 2; start < 2 + count * 3; start += 3) {
        cursors.push({ id: bytes[start], x: bytes[start + 1], y: bytes[start + 2] });
      }
      return { type: "cursors", cursors };
    }
    // 当前服务在局部更新之后发完整棋盘，用完整帧统一更新棋盘、分数与旗数。
    case 0x32:
      return { type: "won" };
    case 0x33:
      requireBytes(4);
      return { type: "lost", id: bytes[1], x: bytes[2], y: bytes[3] };
    case 0xff:
      return { type: "error", message: textDecoder.decode(bytes.subarray(1)) };
    default:
      return null;
  }
}
