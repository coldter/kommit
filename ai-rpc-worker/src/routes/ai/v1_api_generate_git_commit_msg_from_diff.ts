import { isPublicAccess } from "@/middleware/guard";
import { getCommitMsgFromDiff } from "@/pkg/ai/git-msg-generator";
import { errorResponses, successWithDataSchema } from "@/pkg/common/common-responses";
import { createRouteConfig } from "@/pkg/common/route-config";
import type { App } from "@/pkg/hono/app";
import { getCtxLogger } from "@/pkg/lib/context";
import { z } from "@hono/zod-openapi";
import { createWorkersAI } from "workers-ai-provider";

const generateGitCommitMsgFromDiff200ResponseSchema = z.object({});

const route = createRouteConfig({
  tags: ["ai"],
  summary: "Generate git commit msg from diff",
  method: "post",
  path: "/v1/ai.generateGitCommitMsgFromDiff",
  guard: isPublicAccess,
  operationId: "generateGitCommitMsgFromDiff",
  request: {
    body: {
      required: true,
      content: {
        "text/plain": {
          schema: z.string(),
        },
      },
    },
  },
  responses: {
    200: {
      description: "",
      content: {
        "application/json": {
          schema: successWithDataSchema(generateGitCommitMsgFromDiff200ResponseSchema),
        },
      },
    },
    ...errorResponses,
  },
});

export const registerV1ApiGenerateGitCommitMsgFromDiff = (app: App) => {
  app.openapi(route, async (c) => {
    const reqBody = await c.req.text();

    const workersai = createWorkersAI({ binding: c.env.AI });

    const result = await getCommitMsgFromDiff(
      reqBody,
      workersai("@cf/meta/llama-3.1-8b-instruct"),
      getCtxLogger(),
    );

    return c.json({ success: true, data: result }, 200);
  });
};
