import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_PARSER_API_URL: process.env.NEXT_PUBLIC_PARSER_API_URL,
    NEXT_PUBLIC_AGGREGATOR_API_URL: process.env.NEXT_PUBLIC_AGGREGATOR_API_URL,
    NEXT_PUBLIC_WEBSOCKET_URL: process.env.NEXT_PUBLIC_WEBSOCKET_URL,
  },
};

export default nextConfig;
