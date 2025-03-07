import type { Readable } from "node:stream";
import { ReadableStream } from "node:stream/web";

export const createStreamBody = (stream: Readable) => {
  const body = new ReadableStream({
    start(controller) {
      stream.on("data", (chunk) => {
        controller.enqueue(chunk);
      });
      stream.on("end", () => {
        controller.close();
      });
    },

    cancel() {
      stream.destroy();
    },
  });
  return body;
};
