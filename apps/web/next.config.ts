import type { NextConfig } from "next";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.dirname(fileURLToPath(import.meta.url));

const nextConfig: NextConfig = {
  transpilePackages: ["@ai-finops/api-schemas"],
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  ...(process.env.NEXT_STANDALONE_OUTPUT === "1"
    ? {
        output: "standalone" as const,
        outputFileTracingRoot: path.join(rootDir, "../.."),
      }
    : {}),
};

export default nextConfig;
