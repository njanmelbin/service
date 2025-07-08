import { useState } from 'react'
import { Save, RefreshCw, Database, Shield, Bell } from 'lucide-react'
import toast from 'react-hot-toast'

export default function SettingsPage() {
  const [settings, setSettings] = useState({
    database: {
      host: 'localhost',
      port: '5432',
      name: 'postgres',
      maxConnections: '100'
    },
    auth: {
      tokenExpiry: '24',
      requireEmailVerification: true,
      allowPasswordReset: true
    },
    notifications: {
      emailAlerts: true,
      slackIntegration: false,
      webhookUrl: ''
    }
  })

  const handleSave = () => {
    // In a real app, you would save these settings to your backend
    toast.success('Settings saved successfully')
  }

  const handleReset = () => {
    // Reset to default values
    setSettings({
      database: {
        host: 'localhost',
        port: '5432',
        name: 'postgres',
        maxConnections: '100'
      },
      auth: {
        tokenExpiry: '24',
        requireEmailVerification: true,
        allowPasswordReset: true
      },
      notifications: {
        emailAlerts: true,
        slackIntegration: false,
        webhookUrl: ''
      }
    })
    toast.success('Settings reset to defaults')
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
          <p className="mt-1 text-sm text-gray-600">
            Configure your microservices environment
          </p>
        </div>
        <div className="flex space-x-3">
          <button onClick={handleReset} className="btn-secondary flex items-center">
            <RefreshCw className="h-4 w-4 mr-2" />
            Reset
          </button>
          <button onClick={handleSave} className="btn-primary flex items-center">
            <Save className="h-4 w-4 mr-2" />
            Save Changes
          </button>
        </div>
      </div>

      {/* Database Settings */}
      <div className="card">
        <div className="flex items-center mb-4">
          <Database className="h-5 w-5 text-gray-400 mr-2" />
          <h2 className="text-lg font-medium text-gray-900">Database Configuration</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Host</label>
            <input
              type="text"
              value={settings.database.host}
              onChange={(e) => setSettings({
                ...settings,
                database: { ...settings.database, host: e.target.value }
              })}
              className="input-field mt-1"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Port</label>
            <input
              type="text"
              value={settings.database.port}
              onChange={(e) => setSettings({
                ...settings,
                database: { ...settings.database, port: e.target.value }
              })}
              className="input-field mt-1"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Database Name</label>
            <input
              type="text"
              value={settings.database.name}
              onChange={(e) => setSettings({
                ...settings,
                database: { ...settings.database, name: e.target.value }
              })}
              className="input-field mt-1"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Max Connections</label>
            <input
              type="number"
              value={settings.database.maxConnections}
              onChange={(e) => setSettings({
                ...settings,
                database: { ...settings.database, maxConnections: e.target.value }
              })}
              className="input-field mt-1"
            />
          </div>
        </div>
      </div>

      {/* Authentication Settings */}
      <div className="card">
        <div className="flex items-center mb-4">
          <Shield className="h-5 w-5 text-gray-400 mr-2" />
          <h2 className="text-lg font-medium text-gray-900">Authentication</h2>
        </div>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Token Expiry (hours)</label>
            <input
              type="number"
              value={settings.auth.tokenExpiry}
              onChange={(e) => setSettings({
                ...settings,
                auth: { ...settings.auth, tokenExpiry: e.target.value }
              })}
              className="input-field mt-1 max-w-xs"
            />
          </div>
          <div className="space-y-3">
            <div className="flex items-center">
              <input
                type="checkbox"
                checked={settings.auth.requireEmailVerification}
                onChange={(e) => setSettings({
                  ...settings,
                  auth: { ...settings.auth, requireEmailVerification: e.target.checked }
                })}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label className="ml-2 block text-sm text-gray-900">
                Require email verification
              </label>
            </div>
            <div className="flex items-center">
              <input
                type="checkbox"
                checked={settings.auth.allowPasswordReset}
                onChange={(e) => setSettings({
                  ...settings,
                  auth: { ...settings.auth, allowPasswordReset: e.target.checked }
                })}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label className="ml-2 block text-sm text-gray-900">
                Allow password reset
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Notification Settings */}
      <div className="card">
        <div className="flex items-center mb-4">
          <Bell className="h-5 w-5 text-gray-400 mr-2" />
          <h2 className="text-lg font-medium text-gray-900">Notifications</h2>
        </div>
        <div className="space-y-4">
          <div className="space-y-3">
            <div className="flex items-center">
              <input
                type="checkbox"
                checked={settings.notifications.emailAlerts}
                onChange={(e) => setSettings({
                  ...settings,
                  notifications: { ...settings.notifications, emailAlerts: e.target.checked }
                })}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label className="ml-2 block text-sm text-gray-900">
                Email alerts
              </label>
            </div>
            <div className="flex items-center">
              <input
                type="checkbox"
                checked={settings.notifications.slackIntegration}
                onChange={(e) => setSettings({
                  ...settings,
                  notifications: { ...settings.notifications, slackIntegration: e.target.checked }
                })}
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label className="ml-2 block text-sm text-gray-900">
                Slack integration
              </label>
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Webhook URL</label>
            <input
              type="url"
              value={settings.notifications.webhookUrl}
              onChange={(e) => setSettings({
                ...settings,
                notifications: { ...settings.notifications, webhookUrl: e.target.value }
              })}
              placeholder="https://hooks.slack.com/services/..."
              className="input-field mt-1"
            />
          </div>
        </div>
      </div>
    </div>
  )
}