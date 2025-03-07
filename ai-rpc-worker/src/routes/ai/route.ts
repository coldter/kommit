import type { App } from "@/pkg/hono/app";
import { registerV1ApiGenerateGitCommitMsgFromDiff } from "@/routes/ai/v1_api_generate_git_commit_msg_from_diff";

export const setupAiApiRoutes = (app: App) => {
  registerV1ApiGenerateGitCommitMsgFromDiff(app);
};
