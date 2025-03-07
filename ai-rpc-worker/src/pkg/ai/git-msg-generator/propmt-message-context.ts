import type { CoreMessage } from "ai";
import dedent from "dedent";

export const gitDiffToCommitMsgPromptMessagePreContext: CoreMessage[] = [
  {
    role: "user",
    content: dedent`
  File: /tests/userService.test.js
--- /tests/userService.test.js
+++ /tests/userService.test.js
@@ -0,0 +1,20 @@
+const userService = require('../src/userService');
+const assert = require('assert');
+
+describe('User Service', () => {
+  it('should create a new user', () => {
+    const user = userService.createUser('John Doe', 'john@example.com');
+    assert.strictEqual(user.name, 'John Doe');
+    assert.strictEqual(user.email, 'john@example.com');
+  });
+
+  it('should find a user by email', () => {
+    const user = userService.createUser('Jane Doe', 'jane@example.com');
+    const foundUser = userService.findUserByEmail('jane@example.com');
+    assert.strictEqual(foundUser.name, 'Jane Doe');
+  });
+});
  `,
  },

  {
    role: "assistant",
    content: dedent`
    <msg>
test: added unit tests for user service 
</msg>

<body>
- Added user service test for creating new user flow
</body>
    `,
  },
  {
    role: "user",
    content: dedent`
      diff --git a/.gitignore b/.gitignore
index d1e3f8b..a3f1c8a 100644
--- a/.gitignore
+++ b/.gitignore
@@ -1,3 +1,4 @@
 # Ignore node_modules
 node_modules/
 
+# Ignore log files
+*.log"},
    { role: "assistant", content: "<msg>
chore: updated .gitignore
</msg>

<body>
- Added ignore pattern for log files
- Updated ignored files to include log files
</body>"},
    { role: "user", content: "diff --git a/README.md b/README.md
index 83fcf3d..3e2a1b3 100644
--- a/README.md
+++ b/README.md
@@ -10,7 +10,7 @@
 
 ## Installation
 
-Clone the repository and run the install scrpt:
+Clone the repository and run the install script:
 
 \`\`\`bash
 git clone https://github.com/example/repo.git`,
  },
  {
    role: "assistant",
    content: dedent`
    <msg>
docs: updated installation instructions in README
</msg>

<body>
- Updated installation instructions to fix typo in script name
</body>`,
  },
];

export const gitMsgFromDiffSystemPrompt = dedent`
You analyze git patch diffs and generate concise, informative commit messages.

## Input Format
Git patch diff showing file changes with added (+) and removed (-) lines.

## Output Format
Follow the Conventional Commits format:
\`\`\`
<type>[optional scope]: <description>
\`\`\`

## Guidelines
1. **Select appropriate type**: feat, fix, docs, style, refactor, perf, test, chore
2. **Identify scope**: affected component/module according to the project
3. **Write concise description**: imperative mood, no period, under 72 chars
4. **Add details in body** if necessary but prefix with body

## Important
- prefix commit message with \`msg:\`
- prefix detail body with \`body:\`
- only include these two things in response avoid adding any explanation or boiler to response
- there will be last 10 commit message provided if available use that for how to write commit msg it not provided use commit convection
- if previous commit not provided and everything is addition it might be first commit show write commit message like \`chore: project init\` and write every short description in body
- if provided input is other than git diff return empty string in both msg and body tags

## Analysis Tips
- Focus on what changed and why
- Identify patterns across changed files
- Consider function/method names that were modified
- Prioritize significant changes in large diffs

Generate clear, helpful commit messages that accurately represent the changes.

## Example
Input:
\`\`\`diff
File: a/src/components/UserAvatar.js
--- a/src/components/UserAvatar.js
+++ a/src/components/UserAvatar.js
@@ -15,7 +15,7 @@ const UserAvatar = ({ user, size }) => {
   return (
     <Avatar 
       alt={user.name}
-      src={user.avatar}
+      src={user.avatar || '/default-avatar.png'}
       sx={{ width: size, height: size }}
     />
   );
\`\`\`

Output:
\`\`\`
<msg>
fix(components): add fallback for missing user avatars
</msg>

<body>
Added default image '/default-avatar.png' to Avatar component in UserAvatar
</body>
\`\`\`
`;
