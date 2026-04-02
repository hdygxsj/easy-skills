import { useState, useEffect } from 'react'

interface Package {
  id: string
  name: string
  target: string
  source?: string
}

interface Component {
  id: string
  name: string
  type: 'skill' | 'agent' | 'hook' | 'rule'
  packageName?: string
  installed?: boolean
}

type GroupedComponents = {
  skills: Component[]
  agents: Component[]
  hooks: Component[]
  rules: Component[]
}

function App() {
  const [target, setTarget] = useState<'qoder' | 'cursor'>('qoder')
  const [view, setView] = useState<'packages' | 'skills'>('packages')
  const [packages, setPackages] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchPackages()
  }, [target])

  const fetchPackages = async () => {
    setLoading(true)
    // TODO: Replace with actual API call
    // For now, show sample data with type classification
    setPackages([
      {
        name: 'superpowers',
        target: target,
        version: 'v1.0.0',
        components: [
          { id: '1', name: 'brainstorming', type: 'skill', packageName: 'superpowers' },
          { id: '2', name: 'writing-plans', type: 'skill', packageName: 'superpowers' },
          { id: '3', name: 'test-driven-development', type: 'skill', packageName: 'superpowers' },
          { id: '4', name: 'systematic-debugging', type: 'skill', packageName: 'superpowers' },
          { id: '5', name: 'verification-before-completion', type: 'skill', packageName: 'superpowers' },
          { id: '6', name: 'subagent-driven-dev', type: 'agent', packageName: 'superpowers' },
          { id: '7', name: 'team-driven-development', type: 'agent', packageName: 'superpowers' },
          { id: '8', name: 'code-quality-checker', type: 'agent', packageName: 'superpowers' },
        ]
      },
      {
        name: 'open-spec',
        target: target,
        version: 'v2.0.0',
        components: [
          { id: '9', name: 'api-design', type: 'skill', packageName: 'open-spec' },
          { id: '10', name: 'db-schema', type: 'skill', packageName: 'open-spec' },
          { id: '11', name: 'test-gen', type: 'skill', packageName: 'open-spec' },
          { id: '12', name: 'typescript-rules', type: 'rule', packageName: 'open-spec' },
        ]
      }
    ])
    setLoading(false)
  }

  // Group components by type
  const groupComponentsByType = (components: Component[]): GroupedComponents => {
    return {
      skills: components.filter(c => c.type === 'skill'),
      agents: components.filter(c => c.type === 'agent'),
      hooks: components.filter(c => c.type === 'hook'),
      rules: components.filter(c => c.type === 'rule'),
    }
  }

  // Render component group
  const renderGroup = (title: string, icon: string, items: Component[]) => {
    if (items.length === 0) return null
    return (
      <div className="mb-4">
        <h3 className="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-2">
          <span>{icon}</span>
          <span>{title}</span>
          <span className="text-gray-400">({items.length})</span>
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 ml-4">
          {items.map((comp) => (
            <div key={comp.id} className="bg-gray-50 rounded-lg p-3 border border-gray-200 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-green-500"></span>
                <span className="text-sm font-medium">{comp.name}</span>
              </div>
              <span className="text-xs text-gray-400">{comp.packageName}</span>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">Easy Skills Hub</h1>
          <select 
            value={target}
            onChange={(e) => setTarget(e.target.value as 'qoder' | 'cursor')}
            className="px-4 py-2 border border-gray-300 rounded-lg"
          >
            <option value="qoder">Qoder</option>
            <option value="cursor">Cursor</option>
          </select>
        </div>
      </header>

      {/* Install Guide Banner */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg m-6 p-4">
        <h2 className="text-lg font-semibold text-blue-900 mb-2">📦 Install Easy Skills in Your AI IDE</h2>
        <div className="space-y-2 text-sm">
          <div>
            <span className="font-medium text-blue-700">Qoder:</span>
            <code className="ml-2 px-2 py-1 bg-blue-100 rounded text-blue-800">
               Fetch and follow instructions from http://localhost:27842/qoder/easy-skills.md
            </code>
          </div>
          <div>
            <span className="font-medium text-blue-700">Cursor:</span>
            <code className="ml-2 px-2 py-1 bg-blue-100 rounded text-blue-800">
              Fetch and follow instructions from http://localhost:27842/cursor/easy-skills.md
            </code>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="px-6 border-b border-gray-200">
        <nav className="flex space-x-4">
          <button
            onClick={() => setView('packages')}
            className={`py-3 px-1 border-b-2 font-medium ${
              view === 'packages'
                ? 'border-primary text-primary'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Packages
          </button>
          <button
            onClick={() => setView('skills')}
            className={`py-3 px-1 border-b-2 font-medium ${
              view === 'skills'
                ? 'border-primary text-primary'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Skills
          </button>
        </nav>
      </div>

      {/* Content */}
      <main className="p-6">
        {loading ? (
          <div className="text-center py-12 text-gray-500">Loading...</div>
        ) : view === 'packages' ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {packages.map((pkg) => (
              <div key={pkg.name} className="bg-gray-50 rounded-lg p-4 border border-gray-200">
                <h3 className="font-semibold text-lg">{pkg.name}</h3>
                <p className="text-sm text-gray-500 mt-1">{pkg.version}</p>
                <p className="text-sm text-gray-600 mt-2">{pkg.components.length} components</p>
                <div className="mt-3 flex gap-2">
                  <span className="px-2 py-1 bg-green-100 text-green-700 text-xs rounded">
                    ✓ Installed
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="space-y-6">
            {packages.map((pkg: any) => {
              const grouped = groupComponentsByType(pkg.components)
              return (
                <div key={pkg.name} className="border border-gray-200 rounded-lg overflow-hidden">
                  <div className="bg-gray-100 px-4 py-3 border-b border-gray-200 flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-gray-800">{pkg.name}</span>
                      <span className="text-xs text-gray-500">v{pkg.version}</span>
                    </div>
                    <span className="text-sm text-gray-500">{pkg.components.length} components</span>
                  </div>
                  <div className="p-4">
                    {renderGroup('🎯 Skills', '🎯', grouped.skills)}
                    {renderGroup('🤖 Agents', '🤖', grouped.agents)}
                    {renderGroup('🪝 Hooks', '🪝', grouped.hooks)}
                    {renderGroup('📐 Rules', '📐', grouped.rules)}
                    {pkg.components.length === 0 && (
                      <p className="text-gray-400 text-sm">No components</p>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </main>
    </div>
  )
}

export default App
