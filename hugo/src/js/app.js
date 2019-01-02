require('./polyfills')
require('./stimulus')

if (process.env.NODE_ENV === 'production') {
  require('./turbolinks')
}

require('./highlight')
require('./icons')
