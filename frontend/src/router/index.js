import { createBrowserRouter } from 'react-router-dom'
import MainLayout from '../components/MainLayout'
import LoginPage from '../pages/LoginPage'
import HomePage from '../pages/HomePage'
import PatientsPage from '../pages/PatientsPage'
import ExaminationsPage from '../pages/ExaminationsPage'
import ExaminationDetailPage from '../pages/ExaminationDetailPage'
import ReportPage from '../pages/ReportPage'
import ReportsPage from '../pages/ReportsPage'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'patients',
        element: <PatientsPage />,
      },
      {
        path: 'examinations',
        element: <ExaminationsPage />,
      },
      {
        path: 'examinations/:id',
        element: <ExaminationDetailPage />,
      },
      {
        path: 'reports',
        element: <ReportsPage />,
      },
      {
        path: 'reports/:id',
        element: <ReportPage />,
      },
    ],
  },
])
