/** @type {import('next').NextConfig} */
const nextConfig = {
    images: {
        remotePatterns: [
            {
                protocol: 'http',
                hostname: 'localhost',
                port: '9000',
                pathname: '/bluthinator/**'
            },
            {
                protocol: 'http',
                hostname: 'homelab',
                port: '9000',
                pathname: '/bluthinator/**'
            },
            {
                protocol: 'http',
                hostname: 'homelab',
                port: '8000',
                pathname: '/caption/**'
            }
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
