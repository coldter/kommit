import type { App } from "@/pkg/hono/app";
import { setupAiApiRoutes } from "@/routes/ai/route";

export const setupApiRoutes = (app: App) => {
  setupAiApiRoutes(app);
};
