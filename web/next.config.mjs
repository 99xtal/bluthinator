/** @type {import('next').NextConfig} */
const nextConfig = {
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'api.bluthinator.com',
                pathname: '/caption/**'
            },
        ]
    },
    async rewrites() {
        return [
          {
            source: '/img/:episode/:timestamp/:size.jpg',
            destination: `${process.env.IMG_HOST}/frames/:episode/:timestamp/:size.jpg`,
          },
        ];
      },
};

export default nextConfig;
