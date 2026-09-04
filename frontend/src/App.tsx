import { useState, useEffect, useCallback, useRef } from 'react';
import './App.css';
import TopBar from './components/TopBar';
import EditorPane from './components/EditorPane';
import PreviewPane from './components/PreviewPane';
import { type ResumeData, calculateATSScore } from './types';

import { normalizeResumeData } from './utils/normalize';

const API_URL = (import.meta.env.VITE_API_URL || 'http://localhost:8080') as string;

// Default resume data structure
const defaultResumeData: ResumeData = {
  name: '',
  title: '',
  email: '',
  phone: '',
  linkedin: '',
  github: '',
  location: '',
  education: [],
  experiences: [],
  projects: [],
  spokenLanguages: [],
  certifications: [],
  skills: {
    languages: [],
    technologies: []
  }
};

function App() {
  const [resumeData, setResumeData] = useState<ResumeData>(defaultResumeData);
  const [serverStatus, setServerStatus] = useState<'pending' | 'online' | 'offline'>('pending');
  const [isExporting, setIsExporting] = useState(false);
  const [exportError, setExportError] = useState<string | null>(null);
  const [atsScore, setAtsScore] = useState({ score: 0, criteria: [] as Array<{label: string; passed: boolean; points: number}> });
  const [showMobilePreview, setShowMobilePreview] = useState(false);
  const healthCheckIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Load resume data from localStorage on initial load
  useEffect(() => {
    const saved = localStorage.getItem('resume-draft-v2');
    
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setResumeData(normalizeResumeData(parsed));
        return;
      } catch (e) {
        console.error('Failed to parse v2 resume data from localStorage', e);
      }
    }

    // Try fallback to v1
    const oldSaved = localStorage.getItem('resume-draft');
    if (oldSaved) {
      try {
        const parsed = JSON.parse(oldSaved);
        setResumeData(normalizeResumeData(parsed));
      } catch (e) {
        console.error('Failed to parse v1 resume data from localStorage', e);
        setResumeData(defaultResumeData);
      }
    }
  }, []);

  const checkServerStatus = useCallback(async () => {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000);
      const response = await fetch(`${API_URL}/health`, {
        method: 'GET',
        signal: controller.signal
      });
      clearTimeout(timeoutId);
      if (response.ok) {
        setServerStatus('online');
      } else {
        setServerStatus('offline');
      }
    } catch (err: any) {
      if (err.name === 'AbortError') {
        setServerStatus('offline');
        return;
      }
      setServerStatus('offline');
    }
  }, []);

  useEffect(() => {
    // Initial health check
    checkServerStatus();

    // Set up interval for health checks when offline/pending
    healthCheckIntervalRef.current = setInterval(() => {
      checkServerStatus();
    }, 10000);

    return () => {
      if (healthCheckIntervalRef.current) {
        clearInterval(healthCheckIntervalRef.current);
      }
    };
  }, [checkServerStatus]);

  // Save resume data to localStorage whenever it changes (debounced)
  useEffect(() => {
    const handler = setTimeout(() => {
      localStorage.setItem('resume-draft-v2', JSON.stringify(resumeData));
    }, 500);
    return () => clearTimeout(handler);
  }, [resumeData]);

  // Update ATS score whenever resume data changes
  useEffect(() => {
    const result = calculateATSScore(resumeData);
    setAtsScore(result);
  }, [resumeData]);

  // Handle input change for simple fields
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setResumeData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Handle adding an entry to a list field
  const handleAddItem = <T,>(fieldName: keyof ResumeData, newItem: T) => {
    setResumeData((prev) => {
      const currentArray = prev[fieldName] as T[];
      // Enforce max 8 items for list fields
      if (currentArray.length >= 8) {
        return prev;
      }
      return {
        ...prev,
        [fieldName]: [...currentArray, newItem]
      };
    });
  };

  // Handle removing an item from a list field
  const handleRemoveItem = <T,>(fieldName: keyof ResumeData, index: number) => {
    setResumeData((prev) => {
      const currentArray = prev[fieldName] as T[];
      return {
        ...prev,
        [fieldName]: currentArray.filter((_, i) => i !== index)
      };
    });
  };

  // Handle updating a specific field in a list item
  const handleUpdateItemField = <T extends object>(fieldName: keyof ResumeData, index: number, fieldKey: keyof T, value: any) => {
    setResumeData((prev) => {
      if (fieldName === 'skills') {
        return {
          ...prev,
          skills: {
            ...prev.skills,
            [fieldKey as any]: value
          }
        };
      }
      const currentArray = [...(prev[fieldName] as T[])];
      if (index < 0 || index >= currentArray.length) return prev;
      const item = { ...currentArray[index], [fieldKey]: value };
      const newArray = [...currentArray];
      newArray[index] = item;
      return {
        ...prev,
        [fieldName]: newArray as any
      };
    });
  };

  // Handle updating a nested list (e.g., bullets)
  const handleUpdateNestedList = (
    outerFieldName: keyof ResumeData,
    outerIndex: number,
    innerFieldName: keyof any,
    innerIndex: number,
    value: string
  ) => {
    setResumeData(prev => {
      const outerArray = [...(prev[outerFieldName] as any[])];
      if (outerIndex < 0 || outerIndex >= outerArray.length) return prev;
      const innerArray = [...(outerArray[outerIndex][innerFieldName] as string[])];
      if (innerIndex < 0 || innerIndex >= innerArray.length) return prev;
      const newInnerArray = [...innerArray];
      newInnerArray[innerIndex] = value;
      const newOuterItem = { ...outerArray[outerIndex], [innerFieldName]: newInnerArray };
      const newOuterArray = [...outerArray];
      newOuterArray[outerIndex] = newOuterItem;
      return {
        ...prev,
        [outerFieldName]: newOuterArray
      };
    });
  };

  // Handle removing a nested item (e.g., bullet)
  const handleRemoveNestedItem = (
    outerFieldName: keyof ResumeData,
    outerIndex: number,
    innerFieldName: keyof any,
    innerIndex: number
  ) => {
    setResumeData(prev => {
      const outerArray = [...(prev[outerFieldName] as any[])];
      if (outerIndex < 0 || outerIndex >= outerArray.length) return prev;
      const innerArray = [...(outerArray[outerIndex][innerFieldName] as string[])];
      if (innerIndex < 0 || innerIndex >= innerArray.length) return prev;
      const newInnerArray = innerArray.filter((_, i) => i !== innerIndex);
      const newOuterItem = { ...outerArray[outerIndex], [innerFieldName]: newInnerArray };
      const newOuterArray = [...outerArray];
      newOuterArray[outerIndex] = newOuterItem;
      return {
        ...prev,
        [outerFieldName]: newOuterArray
      };
    });
  };

  // Handle updating an array field (e.g., certifications)
  const handleUpdateArrayField = (fieldName: keyof ResumeData, value: any[]) => {
    setResumeData(prev => ({
      ...prev,
      [fieldName]: value
    }));
  };

  // Handle PDF export
  const handleExportPdf = async () => {
    // Disable export if name or title is empty
    if (!resumeData.name.trim() || !resumeData.title.trim()) {
      setExportError('Nome e título são obrigatórios para exportar o PDF.');
      return;
    }

    setIsExporting(true);
    setExportError(null);

    try {
      const response = await fetch(`${API_URL}/generate-pdf`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(resumeData)
      });

      if (!response.ok) {
        const isJson = response.headers.get('content-type')?.includes('application/json');
        if (!isJson) {
          throw new Error('Erro no servidor ao gerar PDF. Tente novamente mais tarde.');
        }
        const errorData = await response.json();
        throw new Error(errorData.message || 'Falha ao gerar PDF');
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'curriculo.pdf';
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (err: any) {
      setExportError(err.message || 'Erro desconhecido ao gerar PDF');
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <div className="app">
      <TopBar
        serverStatus={serverStatus}
        atsScore={atsScore}
        isExporting={isExporting}
        exportError={exportError}
        onExportPdf={handleExportPdf}
        canExport={resumeData.name.trim() !== '' && resumeData.title.trim() !== ''}
      />
      <div className="layout">
        <EditorPane
          resumeData={resumeData}
          onChange={handleChange}
          onAddItem={handleAddItem}
          onRemoveItem={handleRemoveItem}
          onUpdateItemField={handleUpdateItemField}
          onUpdateNestedList={handleUpdateNestedList}
          onRemoveNestedItem={handleRemoveNestedItem}
          onUpdateArrayField={handleUpdateArrayField}
        />
        <PreviewPane
          resumeData={resumeData}
          isMobileOpen={showMobilePreview}
          onCloseMobile={() => setShowMobilePreview(false)}
          onExportPdf={handleExportPdf}
          isExporting={isExporting}
        />
        <button className="fab-preview" onClick={() => setShowMobilePreview(true)}>
          Ver Preview
        </button>
      </div>
    </div>
  );
}

export default App;