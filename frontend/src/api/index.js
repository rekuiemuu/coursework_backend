import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
}

export const patientsAPI = {
  list: (params) => api.get('/patients', { params }),
  get: (id) => api.get(`/patients/${id}`),
  create: (data) => api.post('/patients', data),
  update: (id, data) => api.put(`/patients/${id}`, data),
  delete: (id) => api.delete(`/patients/${id}`),
}

export const examinationsAPI = {
  list: (params) => api.get('/examinations', { params }),
  get: (id) => api.get(`/examinations/${id}`),
  create: (data) => api.post('/examinations', data),
  attachPhotos: (id, photos) => api.post(`/examinations/${id}/photos`, { photos }),
  startAnalysis: (id) => api.post(`/examinations/${id}/analyze`),
  getByPatient: (patientId) => api.get(`/examinations/patient/${patientId}`),
}

export const reportsAPI = {
  list: (params) => api.get('/reports', { params }),
  get: (id) => api.get(`/reports/${id}`),
  getByExamination: (examinationId) => api.get(`/reports/examination/${examinationId}`),
  create: (data) => api.post('/reports', data),
  update: (id, data) => api.put(`/reports/${id}`, data),
}

export default api
