import { useState, useEffect, useRef } from 'preact/hooks'
import { Clipboard, Check, HelpCircle } from 'lucide-preact'
import { EditorView, basicSetup } from 'codemirror'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import './app.css'

const CopyIcon = () => (
  <Clipboard size={12} />
)

const CheckIcon = () => (
  <Check size={12} />
)

const HelpIcon = ({ text }) => (
  <span className="help-icon" data-tooltip={text}>
    <HelpCircle size={14} />
  </span>
)

function CollapsibleSection({ title, children, expanded, onToggle, subSection = false }) {
  return (
    <div className={subSection ? 'subsection-divider' : ''}>
      <div 
        className={subSection ? 'config-subsection-header' : 'config-header'}
        onClick={onToggle}
      >
        {subSection ? (
          <>
            <span className={`chevron-transition ${expanded ? 'rotated' : ''} config-subsection-chevron`}>▶</span>
            <span>{title}</span>
          </>
        ) : (
          <>
            <span>{title}</span>
            <span className={`chevron-transition ${expanded ? 'rotated' : ''}`}>▶</span>
          </>
        )}
      </div>
      
      <div className={`collapsible-grid ${expanded ? 'expanded' : ''}`}>
        <div className="collapsible-content">
          {children}
        </div>
      </div>
    </div>
  )
}

const defaultJSON = [
  {
    "created": "2020-10-03T15:04:05Z",
    "name": "water",
    "description": "quite common on earth",
    "type": "liquid",
    "boiling_point": {
      "units": "C",
      "value": 100
    },
    "secret": null,
    "properties": {
      "tasteless": true,
      "odorless": true,
      "abundant": true,
      "solid": false
    }
  },
  {
    "created": "2020-10-03T15:05:02Z",
    "name": "oxygen",
    "type": "gas",
    "density": {
      "units": "g/L",
      "value": 1.429
    },
    "properties": {
      "abundant": true,
      "odorless": true,
      "solid": false
    }
  },
  {
    "created": "2020-10-03T17:12:37Z",
    "name": "carbon monoxide",
    "type": "gas",
    "dangerous": true,
    "boiling_point": {
      "units": "C",
      "value": -191.5
    },
    "density": {
      "units": "kg/m3",
      "value": 789
    },
    "secret": null,
    "properties": {
      "abundant": true,
      "flammable": true
    }
  }
]

export function App() {
  const [output, setOutput] = useState('Loading...')
  const [optionsExpanded, setOptionsExpanded] = useState(false)
  const [typeExtractionExpanded, setTypeExtractionExpanded] = useState(false)
  const [url, setUrl] = useState('')
  const [loadError, setLoadError] = useState(false)
  const [copied, setCopied] = useState(false)
  
  // Options state
  const [options, setOptions] = useState({
    rootName: 'Document',
    extractCommonTypes: true,
    stringPointersWhenKeyMissing: true,
    skipEmptyKeys: true,
    useMaps: true,
    useMapsMinAttrs: 5,
    timeAsStr: false,
    similarityThreshold: 0.7,
    minSubsetSize: 2,
    minSubsetOccurrences: 2,
    minAddedFields: 2,
  })

  const optionsRef = useRef(options)
  optionsRef.current = options

  const editorRef = useRef(null)
  const editorViewRef = useRef(null)

  // Initialize CodeMirror
  useEffect(() => {
    if (!editorRef.current) return

    const view = new EditorView({
      doc: JSON.stringify(defaultJSON, null, 2),
      extensions: [
        basicSetup,
        json(),
        oneDark,
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            parseJSON()
          }
        })
      ],
      parent: editorRef.current
    })

    editorViewRef.current = view

    return () => view.destroy()
  }, [])

  // Load WASM and parse JSON
  useEffect(() => {
    if (!window.Go) {
      console.error('Go WASM runtime not loaded')
      return
    }

    const go = new Go()
    
    WebAssembly.instantiateStreaming(fetch('json2go.wasm'), go.importObject)
      .then((result) => {
        go.run(result.instance)
        parseJSON()
      })
      .catch((err) => {
        console.error('Failed to load WASM:', err)
        setOutput('Error loading WASM module: ' + err.message)
      })
  }, [])

  const parseJSON = () => {
    if (!editorViewRef.current || !window.json2go) {
      setOutput('waiting for valid json ...')
      return
    }

    const jsonText = editorViewRef.current.state.doc.toString()
    const result = window.json2go(jsonText, optionsRef.current)

    if (result) {
      setOutput(result)
    } else {
      setOutput('waiting for valid json ...')
    }
  }

  useEffect(() => {
    parseJSON()
  }, [options])

  const updateOption = (key, value) => {
    setOptions(prev => ({ ...prev, [key]: value }))
  }

  const loadFromURL = async () => {
    if (!url) return

    try {
      const response = await fetch(url)
      const data = await response.json()
      
      if (editorViewRef.current) {
        const formatted = JSON.stringify(data, null, 2)
        editorViewRef.current.dispatch({
          changes: {
            from: 0,
            to: editorViewRef.current.state.doc.length,
            insert: formatted
          }
        })
      }
      
      setLoadError(false)
      parseJSON()
    } catch (err) {
      console.error('Failed to load URL:', err)
      setLoadError(true)
    }
  }

  const copyToClipboard = () => {
    if (!output || output === 'Loading...' || output.startsWith('waiting')) return
    
    navigator.clipboard.writeText(output).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  return (
    <div is-="view" className="container">
      {/* Header */}
      <div is-="view" box-="square" className="header-box">
        <div className="header-content">
          <img src="gopher.svg" alt="Go Gopher" className="gopher-logo" />
          <div className="header-text">
            <h1 is-="h1" className="header-title">json2go web parser</h1>
            <p is-="p" className="header-description">
              Paste json to generate go struct.<br/>
              For source code and more details see <a href="https://github.com/m-zajac/json2go">https://github.com/m-zajac/json2go</a><br/>
              If you use Visual Studio Code, checkout this cool <a href="https://marketplace.visualstudio.com/items?itemName=m-zajac.vsc-json2go" target="blank">vsc-json2go</a> extension!
            </p>
          </div>
        </div>
      </div>

      {/* Options Panel */}
      <div is-="view" box-="square" className="config-box" shear-="top">
        <div class="header" onClick={() => setOptionsExpanded(!optionsExpanded)} style={{ cursor: 'pointer' }}>
          <span is-="badge" variant-="foreground2">
            CONFIGURATION OPTIONS
            <span className={`chevron-transition ${optionsExpanded ? 'rotated' : ''}`} style={{ padding: '0 0.5ch' }}>&gt;</span>
          </span>
        </div>
        <div className={`collapsible-grid ${optionsExpanded ? 'expanded' : ''}`}>
          <div className="collapsible-content">
            <div className="config-content options-columns">
              <div className="options-group">
                <div className="option-row">
                  <label>
                    Root name
                    <HelpIcon text="Name of the generated root Go struct." />
                  </label>
                  <input 
                    is-="input"
                    type="text" 
                    value={options.rootName}
                    onInput={(e) => updateOption('rootName', e.target.value)}
                  />
                </div>

                <div className="option-row">
                  <label>
                    String pointers
                    <HelpIcon text="Toggles whether missing string keys in some documents should result in a pointer (*string) instead of a regular string." />
                  </label>
                  <label className="checkbox-label">
                    <input 
                      type="checkbox"
                      checked={options.stringPointersWhenKeyMissing}
                      onChange={(e) => updateOption('stringPointersWhenKeyMissing', e.target.checked)}
                    />
                    Use *string when missing
                  </label>
                </div>

                <div className="option-row">
                  <label>
                    Skip empty
                    <HelpIcon text="Toggles skipping keys in the input that contained only null values." />
                  </label>
                  <label className="checkbox-label">
                    <input 
                      type="checkbox"
                      checked={options.skipEmptyKeys}
                      onChange={(e) => updateOption('skipEmptyKeys', e.target.checked)}
                    />
                    Skip keys with null values
                  </label>
                </div>
              </div>

              <div className="options-group">
                <div className="option-row">
                  <label>
                    Use maps
                    <HelpIcon text="Defines if the parser should attempt to use maps instead of structs when objects can be represented as map[string]T." />
                  </label>
                  <label className="checkbox-label">
                    <input 
                      type="checkbox"
                      checked={options.useMaps}
                      onChange={(e) => updateOption('useMaps', e.target.checked)}
                    />
                    Try maps for objects
                  </label>
                </div>

                <div className="option-row">
                  <label>
                    Min map attributes
                    <HelpIcon text="Minimum number of attributes an object must have to be considered for conversion to a map." />
                  </label>
                  <input 
                    is-="input"
                    type="number"
                    min="1"
                    step="1"
                    value={options.useMapsMinAttrs}
                    onInput={(e) => updateOption('useMapsMinAttrs', parseInt(e.target.value))}
                  />
                </div>

                <div className="option-row">
                  <label>
                    Time format
                    <HelpIcon text="Toggles whether to treat valid time strings as time.Time or just as strings." />
                  </label>
                  <label className="checkbox-label">
                    <input 
                      type="checkbox"
                      checked={options.timeAsStr}
                      onChange={(e) => updateOption('timeAsStr', e.target.checked)}
                    />
                    Use string for time.Time
                  </label>
                </div>
              </div>
            </div>

            <CollapsibleSection 
              title="Type Extraction Settings" 
              expanded={typeExtractionExpanded} 
              onToggle={() => setTypeExtractionExpanded(!typeExtractionExpanded)}
              subSection
            >
              <div className="config-content options-columns">
                <div className="options-group">
                  <div className="option-row">
                    <label>
                      Common types
                      <HelpIcon text="Toggles extracting common JSON nodes as separate types." />
                    </label>
                    <label className="checkbox-label">
                      <input 
                        type="checkbox"
                        checked={options.extractCommonTypes}
                        onChange={(e) => updateOption('extractCommonTypes', e.target.checked)}
                      />
                      Extract & reuse structs
                    </label>
                  </div>

                  <div className="option-row">
                    <label>
                      Similarity threshold
                      <HelpIcon text="Minimum similarity score (0.0 to 1.0) required to consider two types as similar enough to be merged." />
                    </label>
                    <input 
                      is-="input"
                      type="number"
                      min="0"
                      max="1"
                      step="0.1"
                      value={options.similarityThreshold}
                      onInput={(e) => updateOption('similarityThreshold', parseFloat(e.target.value))}
                    />
                  </div>
                  
                  <div className="option-row">
                    <label>
                      Min subset size
                      <HelpIcon text="Minimum number of fields a shared subset must have to be considered for extraction as a separate type." />
                    </label>
                    <input 
                      is-="input"
                      type="number"
                      min="1"
                      step="1"
                      value={options.minSubsetSize}
                      onInput={(e) => updateOption('minSubsetSize', parseInt(e.target.value))}
                    />
                  </div>
                </div>

                <div className="options-group">
                  <div className="option-row">
                    <label>
                      Min occurrences
                      <HelpIcon text="Minimum number of times a shared subset must occur across different objects to be extracted." />
                    </label>
                    <input 
                      is-="input"
                      type="number"
                      min="1"
                      step="1"
                      value={options.minSubsetOccurrences}
                      onInput={(e) => updateOption('minSubsetOccurrences', parseInt(e.target.value))}
                    />
                  </div>

                  <div className="option-row">
                    <label>
                      Min added fields
                      <HelpIcon text="Minimum number of unique fields a new type must have to justify its extraction." />
                    </label>
                    <input 
                      is-="input"
                      type="number"
                      min="1"
                      step="1"
                      value={options.minAddedFields}
                      onInput={(e) => updateOption('minAddedFields', parseInt(e.target.value))}
                    />
                  </div>
                </div>
              </div>
            </CollapsibleSection>
          </div>
        </div>
      </div>

      {/* Main Content Areas */}
      <div className="main-layout">
        {/* JSON Input */}
        <div is-="view" box-="square" className="panel" shear-="top">
          <div class="header">
              <span is-="badge" variant-="foreground2">JSON INPUT</span>
          </div>
          <div className="panel-toolbar">
            <div className="input-group">
              <input 
                is-="input"
                type="text"
                placeholder="https://api.example.com/data.json"
                value={url}
                onInput={(e) => setUrl(e.target.value)}
                className="input-fullwidth"
              />
              <button 
                is-="button"
                variant-={loadError ? 'red' : 'blue'}
                onClick={loadFromURL}
                className="button-fixed"
                size-="small"
              >
                LOAD FROM URL
              </button>
            </div>
          </div>
          <div ref={editorRef} className="panel-editor" />
        </div>

        {/* Go Output */}
        <div is-="view" box-="square" className="panel" shear-="top">
          <div class="header">
            <span is-="badge" variant-="foreground2">GO STRUCT OUTPUT</span>
          </div>
          
          <div className="output-container">
            <pre is-="pre" className="xoutput-pre">
              <button 
                is-="button" 
                className="copy-button"
                variant-={copied ? 'green' : 'blue'}
                size-="small"
                onClick={copyToClipboard}
                disabled={!output || output === 'Loading...' || output.startsWith('waiting')}
                title={copied ? 'Copied!' : 'Copy to clipboard'}
              >
                {copied ? <CheckIcon /> : <CopyIcon />}
              </button>
              {output}
            </pre>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="footer">
        Version _VERSION_
        <br/>
        Thanks to <a href="https://github.com/MariaLetta/free-gophers-pack">Maria Letta</a> for awesome gopher image! :)
      </div>
    </div>
  )
}
