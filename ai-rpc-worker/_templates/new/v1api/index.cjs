const tsx = require("tsx/cjs/api");

// Register tsx enhancement
const unregister = tsx.register();

const loaded = require("./index.ts");

module.exports = loaded;
// Unregister when needed
unregister();
