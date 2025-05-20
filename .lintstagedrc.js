export default {
  "(services|shared)/**/*.go": (files) => {
    const directories = [
      ...new Set(files.map((file) => file.replace(/\/[^/]+\.go$/, ""))),
    ];

    return directories.flatMap((dir) => [
      `golangci-lint fmt ${dir}`,
      `golangci-lint run --fix ${dir}`,
    ]);
  },
};
