/** @type {import('next').NextConfig} */
const nextConfig = {};

module.exports = nextConfig;

module.exports = {
  serverRuntimeConfig: {
    // Will only be available on the server side
    baseAPI: "backend-service.backend.svc.cluster.local:2022",
  },
};
