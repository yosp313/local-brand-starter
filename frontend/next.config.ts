import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "myaws-default-bucket-1.s3.eu-central-1.amazonaws.com",
      },
    ],
  },
};

export default nextConfig;
