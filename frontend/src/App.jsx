import { Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './components/MainLayout'
import LoginPage from './pages/LoginPage'
import HomePage from './pages/HomePage'
import PatientsPage from './pages/PatientsPage'
import ExaminationsPage from './pages/ExaminationsPage'
import ExaminationDetailPage from './pages/ExaminationDetailPage'
import ReportPage from './pages/ReportPage'

function ProtectedRoute({ children }) {
  const token = localStorage.getItem('token')
  return token ? children : <Navigate to="/login" replace />
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<ProtectedRoute><MainLayout /></ProtectedRoute>}>
        <Route index element={<Navigate to="/examinations" replace />} />
        <Route path="home" element={<HomePage />} />
        <Route path="patients" element={<PatientsPage />} />
        <Route path="examinations" element={<ExaminationsPage />} />
        <Route path="examinations/:id" element={<ExaminationDetailPage />} />
        <Route path="reports/:id" element={<ReportPage />} />
      </Route>
    </Routes>
  )
}

export default App
