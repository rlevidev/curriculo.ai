import { useState, useEffect, useCallback, useRef } from 'react';
import './App.css';
import TopBar from './components/TopBar';
import EditorPane from './components/EditorPane';
import PreviewPane from './components/PreviewPane';
import { type ResumeData, calculateATSScore } from './types';

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
  const [overflowWarning, setOverflowWarning] = useState<string | null>(null);
  const [showMobilePreview, setShowMobilePreview] = useState(false);
  const previewRef = useRef<HTMLDivElement>(null);
  const healthCheckIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const exportTimeoutRef = useRef<ReturnType<typeof setTimeout> | null | undefined>(null);

  // Load resume data from localStorage on initial load
  useEffect(() => {
    const saved = localStorage.getItem('resume-draft');
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        // Validate and merge with default structure to ensure all fields exist
        setResumeData({ ...defaultResumeData, ...parsed });
      } catch (e) {
        console.error('Failed to parse resume data from localStorage', e);
        setResumeData(defaultResumeData);
      }
    }

    // Initial health check
    checkServerStatus();

    // Set up interval for health checks when offline/pending
    healthCheckIntervalRef.current = setInterval(() => {
      if (serverStatus !== 'online') {
        checkServerStatus();
      }
    }, 10000);

    return () => {
      if (healthCheckIntervalRef.current) {
        clearInterval(healthCheckIntervalRef.current);
      }
      if (exportTimeoutRef.current) {
        clearTimeout(exportTimeoutRef.current);
      }
    };
  }, [serverStatus]);

  // Save resume data to localStorage whenever it changes (debounced)
  useEffect(() => {
    const handler = setTimeout(() => {
      localStorage.setItem('resume-draft', JSON.stringify(resumeData));
    }, 500);
    return () => clearTimeout(handler);
  }, [resumeData]);

  // Update ATS score whenever resume data changes
  useEffect(() => {
    const result = calculateATSScore(resumeData);
    setAtsScore(result);
  }, [resumeData]);

  // Check for overflow in preview pane
  useEffect(() => {
    if (!previewRef.current) return;

    const checkOverflow = () => {
      const previewContainer = previewRef.current;
      if (!previewContainer) return;

      // Approximate A4 height at 96dpi: 1123px
      const overflowThreshold = 1123;
      const scrollHeight = previewContainer.scrollHeight;

      if (scrollHeight > overflowThreshold) {
        const estimatedPages = Math.ceil(scrollHeight / overflowThreshold);
        setOverflowWarning(`⚠ Conteúdo pode ocupar ~${estimatedPages} páginas. Considere reduzir o texto.`);
      } else {
        setOverflowWarning(null);
      }
    };

    // Check after DOM updates
    const observer = new ResizeObserver(checkOverflow);
    observer.observe(previewRef.current);

    // Initial check
    checkOverflow();

    return () => observer.disconnect();
  }, [resumeData]); // Re-run when resume data changes (which triggers preview update)

  // Check server status
  const checkServerStatus = useCallback(async () => {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/health`, {
        method: 'GET',
        signal: AbortSignal.timeout(5000)
      });
      if (response.ok) {
        setServerStatus('online');
      } else {
        setServerStatus('offline');
      }
    } catch (err) {
      setServerStatus('offline');
    }
  }, []);

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

    // Change button text after 10 seconds if still exporting
    exportTimeoutRef.current = setTimeout(() => {
      // This will be handled by the button's loading state
    }, 10000);

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/generate-pdf`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(resumeData)
      });

      if (!response.ok) {
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
      if (exportTimeoutRef.current) {
        clearTimeout(exportTimeoutRef.current);
        exportTimeoutRef.current = undefined;
      }
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
          previewRef={previewRef}
          overflowWarning={overflowWarning}
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