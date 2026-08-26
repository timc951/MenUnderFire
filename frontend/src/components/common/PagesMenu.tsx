import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { usePublicSitePages } from '../../hooks/useSitePages';

interface PagesMenuProps {
  className?: string;
}

export function PagesMenu({ className = '' }: PagesMenuProps) {
  const [isOpen, setIsOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const { pages, isLoading } = usePublicSitePages();

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Close menu on escape key
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setIsOpen(false);
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, []);

  const handleNavigate = (path: string) => {
    navigate(path);
    setIsOpen(false);
  };

  return (
    <div ref={menuRef} className={`relative ${className}`}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="p-2 text-stone-300 hover:text-white transition-colors"
        aria-expanded={isOpen}
        aria-haspopup="true"
        aria-label="Menu"
      >
        {/* Hamburger icon */}
        <svg
          className="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d={isOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"}
          />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-2 w-48 bg-stone-900/95 border border-stone-700 rounded-lg shadow-xl backdrop-blur-sm z-50">
          <ul className="py-2">
            <li>
              <button
                onClick={() => handleNavigate('/')}
                className="w-full text-left px-4 py-2 text-stone-300 hover:text-white hover:bg-stone-800 transition-colors"
              >
                Home
              </button>
            </li>
            {isLoading ? (
              <li className="px-4 py-2 text-stone-400 text-sm">Loading...</li>
            ) : (
              pages.map((page) => (
                <li key={page.id}>
                  <button
                    onClick={() => handleNavigate(`/page/${page.slug}`)}
                    className="w-full text-left px-4 py-2 text-stone-300 hover:text-white hover:bg-stone-800 transition-colors"
                  >
                    {page.title}
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      )}
    </div>
  );
}
