import {
  gitDiffToCommitMsgPromptMessagePreContext,
  gitMsgFromDiffSystemPrompt,
} from "@/pkg/ai/git-msg-generator/propmt-message-context";
import { extractFromTag } from "@/pkg/ai/utils/tag";
import { type CoreMessage, type LanguageModel, generateText } from "ai";

function getPromptMessages(diff: string): Array<CoreMessage> {
  return [
    ...gitDiffToCommitMsgPromptMessagePreContext,
    {
      role: "user",
      content: diff,
    },
  ];
}

interface Logger {
  info: (message: string, meta: Record<string, any>) => void;
}
export async function getCommitMsgFromDiff(diff: string, model: LanguageModel, logger: Logger) {
  const generatedTextResult = await generateText({
    model,
    system: gitMsgFromDiffSystemPrompt,
    messages: getPromptMessages(diff),
    maxTokens: 512,
  });

  const generatedText = generatedTextResult.text;
  logger?.info("[getCommitMsgFromDiff] Generated commit message", {
    usage: {
      ...generatedTextResult.usage,
    },
  });

  return {
    msg: extractFromTag("msg", generatedText),
    body: extractFromTag("body", generatedText),
  };
}
