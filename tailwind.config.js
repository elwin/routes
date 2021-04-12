module.exports = {
    purge: [
        'resources/**/*.html',
        'resources/**/*.js',
    ],
    darkMode: false, // or 'media' or 'class'
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter var'],
            },
        },
    },
    variants: {
        extend: {},
    },
    plugins: [],
}
