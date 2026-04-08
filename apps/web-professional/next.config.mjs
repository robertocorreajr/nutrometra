/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: [
    "@nutrometra/ui",
    "@nutrometra/api-client",
    "@nutrometra/auth",
  ],
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8081"}/:path*`,
      },
    ]
  },
}

export default nextConfig
