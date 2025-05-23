export default {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "header-max-length": [2, "always", 60],
    "body-max-line-length": [0, "always"],
  },
};
