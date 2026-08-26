import { ReactNode, useState } from 'react';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { RbacProvider } from '../../rbac';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <RbacProvider>
      <div className="h-screen flex flex-col bg-stone-100 dark:bg-stone-900">
        <Header onToggleSidebar={() => setSidebarOpen(!sidebarOpen)} />
        <div className="flex flex-1 overflow-hidden">
          <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
          <main className="flex-1 p-4 md:p-6 min-w-0 overflow-y-auto">
            {children}
          </main>
        </div>
      </div>
    </RbacProvider>
  );
}
