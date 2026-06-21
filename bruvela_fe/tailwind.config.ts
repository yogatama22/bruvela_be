import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  theme: {
    extend: {
      colors: {
        bruvela: {
          50: '#fbf9f7',
          100: '#f5f0ec',
          200: '#eee7e0',
          300: '#dcd1c8',
          400: '#9c8e82',
          500: '#6d5d52',
          600: '#4a3e35',
          700: '#3d3329',
          800: '#2f271f',
          900: '#221b15',
          950: '#16100c',
        },
      },
    },
  },
}
