import { render } from 'preact'
import { App } from './app.jsx'
import './index.css'

// Set WebTUI theme
document.documentElement.setAttribute('data-webtui-theme', 'catppuccin-mocha')

render(<App />, document.getElementById('app'))
