import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  reactCompiler: true,
  images: {
    domains: ['storage.googleapis.com'],
  },
};

export default nextConfig;
