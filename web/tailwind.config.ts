import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/elements/**/*.{js,ts,jsx,tsx,mdx}"
  ],
  theme: {
    extend: {
      backgroundImage: {
        "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
        "gradient-conic":
          "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
      },
      colors: {
        "theme-white": "#F5F5F5",
        "theme-red": "#C74A3E",
        "theme-red-shadow": "#99372E",
        "theme-orange": "#FEA34b",
        "theme-yellow": "#F8E5AD",
        "theme-black": "#030303"
      },
      rotate: {
        '-45': '-45deg',
        '-30': '-30deg',
        '-15': '-15deg',
        '-10': '-10deg',
        '-5': '-5deg',
        '-2.5': '-2.5deg',
        '0': '0deg',
        '2.5': '2.5deg',
        '5': '5deg',
        '10': '10deg',
        '15': '15deg',
        '30': '30deg',
        '45': '45deg',
      },
    },
  },
  plugins: [],
};
export default config;
