import { useEffect, useState } from 'react'
import { Server, Users, Activity, AlertCircle } from 'lucide-react'
import { healthApi } from '../services/api'

interface ServiceStatus {
  name: string
  status: 'healthy' | 'unhealthy' | 'unknown'
  lastChecked: Date
}

export default function DashboardPage() {
  const [services, setServices] = useState<ServiceStatus[]>([
    { name: 'Sales Service', status: 'unknown', lastChecked: new Date() },
    { name: 'Auth Service', status: 'unknown', lastChecked: new Date() }
  ])

  const checkServiceHealth = async () => {
    try {
      await healthApi.checkLiveness()
      setServices(prev => prev.map(service => 
        service.name === 'Sales Service' 
          ? { ...service, status: 'healthy' as const, lastChecked: new Date() }
          : service
      ))
    } catch (error) {
      setServices(prev => prev.map(service => 
        service.name === 'Sales Service' 
          ? { ...service, status: 'unhealthy' as const, lastChecked: new Date() }
          : service
      ))
    }
  }

  useEffect(() => {
    checkServiceHealth()
    const interval = setInterval(checkServiceHealth, 30000) // Check every 30 seconds
    return () => clearInterval(interval)
  }, [])

  const stats = [
    {
      name: 'Total Services',
      value: '2',
      icon: Server,
      color: 'text-blue-600',
      bgColor: 'bg-blue-100'
    },
    {
      name: 'Active Users',
      value: '12',
      icon: Users,
      color: 'text-green-600',
      bgColor: 'bg-green-100'
    },
    {
      name: 'Healthy Services',
      value: services.filter(s => s.status === 'healthy').length.toString(),
      icon: Activity,
      color: 'text-emerald-600',
      bgColor: 'bg-emerald-100'
    },
    {
      name: 'Alerts',
      value: services.filter(s => s.status === 'unhealthy').length.toString(),
      icon: AlertCircle,
      color: 'text-red-600',
      bgColor: 'bg-red-100'
    }
  ]

  const getStatusColor = (status: ServiceStatus['status']) => {
    switch (status) {
      case 'healthy':
        return 'text-green-600 bg-green-100'
      case 'unhealthy':
        return 'text-red-600 bg-red-100'
      default:
        return 'text-gray-600 bg-gray-100'
    }
  }

  const getStatusText = (status: ServiceStatus['status']) => {
    switch (status) {
      case 'healthy':
        return 'Healthy'
      case 'unhealthy':
        return 'Unhealthy'
      default:
        return 'Unknown'
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-600">
          Overview of your microservices infrastructure
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <div key={stat.name} className="card">
            <div className="flex items-center">
              <div className={`flex-shrink-0 p-3 rounded-lg ${stat.bgColor}`}>
                <stat.icon className={`h-6 w-6 ${stat.color}`} />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">{stat.name}</p>
                <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Service Status */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-gray-900">Service Status</h2>
          <button
            onClick={checkServiceHealth}
            className="btn-secondary text-sm"
          >
            Refresh
          </button>
        </div>
        <div className="space-y-3">
          {services.map((service) => (
            <div key={service.name} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
              <div className="flex items-center">
                <Server className="h-5 w-5 text-gray-400 mr-3" />
                <div>
                  <p className="text-sm font-medium text-gray-900">{service.name}</p>
                  <p className="text-xs text-gray-500">
                    Last checked: {service.lastChecked.toLocaleTimeString()}
                  </p>
                </div>
              </div>
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(service.status)}`}>
                {getStatusText(service.status)}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Recent Activity */}
      <div className="card">
        <h2 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h2>
        <div className="space-y-3">
          <div className="flex items-center text-sm">
            <div className="w-2 h-2 bg-green-400 rounded-full mr-3"></div>
            <span className="text-gray-600">Sales service started successfully</span>
            <span className="ml-auto text-gray-400">2 minutes ago</span>
          </div>
          <div className="flex items-center text-sm">
            <div className="w-2 h-2 bg-blue-400 rounded-full mr-3"></div>
            <span className="text-gray-600">New user registered</span>
            <span className="ml-auto text-gray-400">5 minutes ago</span>
          </div>
          <div className="flex items-center text-sm">
            <div className="w-2 h-2 bg-yellow-400 rounded-full mr-3"></div>
            <span className="text-gray-600">Database migration completed</span>
            <span className="ml-auto text-gray-400">1 hour ago</span>
          </div>
        </div>
      </div>
    </div>
  )
}