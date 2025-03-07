import dedent from "dedent";

export function extractFromTag(tag: string, content: string) {
  const regex = new RegExp(`<${tag}>([\\s\\S]*?)<\/${tag}>`, "i");

  const match = content.match(regex);

  return dedent(match?.[1] ?? "");
}
