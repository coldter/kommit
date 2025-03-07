// ! If call outside of hono request context these will throw context not found error

// import type { HonoEnv } from "@/pkg/hono/env";
// import { getContext } from "hono/context-storage";

/**
 * @throws {Error} context not found
 */
// TODO: remove this or add logger
// export const getCtxLogger = () => {
//   return getContext<HonoEnv>().get("logger");
// };
export const getCtxLogger = () => {
  return console;
};
