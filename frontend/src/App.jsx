import { Routes, Route } from 'react-router-dom'
import MainLayout from './components/MainLayout'
import LoginPage from './pages/LoginPage'
import HomePage from './pages/HomePage'
import PatientsPage from './pages/PatientsPage'
import ExaminationsPage from './pages/ExaminationsPage'
import ExaminationDetailPage from './pages/ExaminationDetailPage'
import ReportPage from './pages/ReportPage'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<MainLayout />}>
        <Route index element={<HomePage />} />
        <Route path="patients" element={<PatientsPage />} />
        <Route path="examinations" element={<ExaminationsPage />} />
        <Route path="examinations/:id" element={<ExaminationDetailPage />} />
        <Route path="reports/:id" element={<ReportPage />} />
      </Route>
    </Routes>
  )
}

export default App
