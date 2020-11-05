//import { init as initApm } from 'elastic-apm-js-base'
import { init as initApm } from '@elastic/apm-rum'

const apm = initApm({
  serviceName: 'frontend-react',
  serviceVersion: '0.90',
  serverUrl: 'http://localhost:8200',
  distributedTracingOrigins: ['http://localhost:8080'],
  debug: true
})

export default apm;
