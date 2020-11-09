import { BACKEND_URL, APM_SERVER_URL } from './constants.js'
import { init as initApm } from '@elastic/apm-rum'

const apm = initApm({
  serviceName: 'frontend-react',
  serviceVersion: '0.90',
  serverUrl: APM_SERVER_URL,
  distributedTracingOrigins: [BACKEND_URL],
  debug: true
})

export default apm;
