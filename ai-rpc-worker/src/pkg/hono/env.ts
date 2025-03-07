export type Env = {
  RATE_LIMIT_CACHE: KVNamespace;
  AI: Ai;
};

export type HonoEnv = {
  Bindings: Env;
  Variables: Record<string, never>;
};
