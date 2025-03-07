import * as changeCase from "change-case";
import type Enquirer from "enquirer";
import fs from "fs-extra";

function listDirectories(directoryPath: string): string[] {
  const entries = fs.readdirSync(directoryPath, { withFileTypes: true });
  return entries.filter((entry) => entry.isDirectory()).map((entry) => entry.name);
}

export default {
  prompt: ({
    inquirer,
  }: {
    inquirer: Enquirer;
  }) => {
    const d = listDirectories("src/routes");
    d.push("new directory");

    const questions = [
      {
        name: "api_name",
        message: "API name",
        type: "input",
      },
      {
        name: "route_dir",
        message: "Route directory",
        type: "select",
        choices: structuredClone(d),
      },
    ] satisfies Parameters<(typeof Enquirer)["prompt"]>[0];

    return inquirer
      .prompt(questions)
      .then(
        (answers: {
          api_name: string;
          route_dir: string;
        }) => {
          if (answers.route_dir !== "new directory") {
            return answers;
          }

          return inquirer
            .prompt({
              name: "new_dir",
              message: "New directory",
              type: "input",
            })
            .then(
              (newDir: {
                new_dir: string;
              }) => {
                return {
                  ...answers,
                  route_dir: newDir.new_dir,
                };
              },
            );
        },
      )
      .then(
        (answers: {
          api_name: string;
          route_dir: string;
        }) => {
          let { api_name, route_dir } = answers;
          // * remove the v1_ prefix
          api_name = api_name.replace(/^v1_api_/, "");
          const apiOperationIdPath = api_name;
          // * convert to snake_case
          api_name = changeCase.snakeCase(api_name);

          if (!d.includes(route_dir)) {
            // * create the directory
            fs.mkdirSync(`src/routes/${route_dir}`);
          }

          const routePath = `./${route_dir}/v1_api_${api_name}.js`;
          const absRouteDirPath = `src/routes/${route_dir}`;
          const absPath = `src/routes/${routePath}`;

          const registerFunctionName = `registerV1Api${changeCase.pascalCase(api_name)}`;
          const setupApiFunctionName = `setup${changeCase.pascalCase(route_dir)}ApiRoutes`;
          const apiFileName = `v1_api_${api_name}.ts`;
          const apiNameSpace = route_dir;
          const apiOperationId = changeCase.camelCase(api_name);
          // const apiOperationIdPath = api_name
          //   .split("_")
          //   .map((part) => changeCase.camelCase(part))
          //   .join(".");

          const absApiFilePath = `src/routes/${route_dir}/${apiFileName}`;

          return {
            ...answers,
            routePath,
            absPath,
            absRouteDirPath,
            registerFunctionName,
            setupApiFunctionName,
            apiFileName,
            apiNameSpace,
            apiOperationId,
            absApiFilePath,
            apiOperationIdPath,
          };
        },
      );
  },
};
