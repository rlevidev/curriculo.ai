import React, { useState, useEffect } from 'react';
import type { ATSResult } from '../types';

interface TopBarProps {
  serverStatus: 'pending' | 'online' | 'offline';
  atsScore: ATSResult;
  isExporting: boolean;
  exportError: string | null;
  onExportPdf: () => void;
}

const TopBar: React.FC<TopBarProps> = ({
  serverStatus,
  atsScore,
  isExporting,
  exportError,
  onExportPdf
}) => {
  const [exportErrorVisible, setExportErrorVisible] = useState(false);

  // Show export error for 5 seconds
  useEffect(() => {
    if (exportError) {
      setExportErrorVisible(true);
      const timer = setTimeout(() => setExportErrorVisible(false), 5000);
      return () => clearTimeout(timer);
    } else {
      setExportErrorVisible(false);
    }
  }, [exportError]);

  return (
    <div className="topbar">
      <div className="brand">
        <span className="dot">◆</span> curriculo<span style={{ color: 'var(--text-faint)' }}>.ai</span>
      </div>
      <nav className="tabs">
        <button className="active">Editor</button>
        <button>Templates</button>
        <button>Histórico</button>
      </nav>
      <div className="topbar-right">
        <div className="compile-status" id="compileStatus">
          <span className={`dot ${serverStatus === 'online' ? 'ok' : serverStatus === 'offline' ? '' : 'busy'}`}>
          </span>
          <span id="compileLabel">
            {serverStatus === 'online' ? 'servidor online' :
              serverStatus === 'offline' ? 'servidor offline' : 'servidor acordando...'}
          </span>
        </div>
        <button
          className="btn-export"
          onClick={onExportPdf}
          disabled={isExporting || !(atsScore.score > 0)} // Disable if exporting or no data
        >
          {isExporting ? 'Exportando PDF...' : 'Exportar PDF ↓'}
        </button>
        {exportErrorVisible && exportError && (
          <div className="export-error" style={{
            position: 'absolute',
            top: '60px',
            right: '20px',
            background: 'var(--amber)',
            color: '#04211c',
            padding: '8px 12px',
            borderRadius: '6px',
            fontFamily: 'var(--font-mono)',
            fontSize: '12px',
            whiteSpace: 'nowrap'
          }}>
            {exportError}
          </div>
        )}
      </div>
    </div>
  );
};

export default TopBar;