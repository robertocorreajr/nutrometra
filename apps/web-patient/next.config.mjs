/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: [
    "@nutrometra/ui",
    "@nutrometra/api-client",
    "@nutrometra/auth",
  ],
}

export default nextConfig
