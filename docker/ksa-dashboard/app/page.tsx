'use client';
import FilterPanel from "./components/FilterPanel";
import React, { useEffect, useState } from "react";
import { Misconfiguration } from "./types/Misconfiguration";
import { Vulnerability } from "./types/Vulnerability";
import MisconfigurationList from "./components/MisconfigurationList";
import VulnerabilityList from "./components/VulnerabilityList";
import { Report } from "./types/Report";
import Toast from "./components/Toast";
import SidePanel from "./components/SidePanel";
import ConfirmationDialog from "./components/ConfirmationDialog";
import Skeleton from 'react-loading-skeleton'
import 'react-loading-skeleton/dist/skeleton.css'
import { Item } from "./types/Item";
import { ITEM_STATUS_FAIL, ITEM_STATUS_PASS, ITEM_STATUS_RESOLVED, ITEM_STATUS_MANUAL, RERUN_SCAN_CONFIRMATION, RESOLVE_ALL_CONFIRMATION, DELETE_ALL_CONFIRMATION, ITEM_TYPE_MISCONFIGURATION, ALL_SCANNERS } from "./constants";

export default function Home() { 
  // Filter Panel
  const [selectedScanner, setSelectedScanner] = useState("Trivy");
  const [selectedType, setSelectedType] = useState(ITEM_TYPE_MISCONFIGURATION);
  const [selectedReport, setSelectedReport] = useState("");
  const [showResolved, setShowResolved] = useState(false);
  const [query, setQuery] = useState("");

  // Data
  const [reports, setReports] = useState<Report[]>([]);
  const [misconfigurations, setMisconfigurations] = useState<Record<string, Misconfiguration[]> | null>(null);
  const [vulnerabilities, setVulnerabilities] = useState<Record<string, Vulnerability[]> | null>(null);

  // Rendering
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [jobInProgress, setJobInProgress] = useState(false);
  const [openCategory, setOpenCategory] = useState<number | null>(null);
  const [openItem, setOpenItem] = useState<number | null>(null);

  // Interactivity
  const [showToast, setShowToast] = useState(false);
  const [message, setMessage] = useState("");
  const [showConfirmationDialog, setShowConfirmationDialog] = useState(false);
  const [confirmationDialogText, setConfirmationDialogText] = useState("");
  const [confirmationDialogAction, setConfirmationDialogAction] = useState<() => void>(() => {});

  const fetchData = async () => {
    if (selectedReport === "" && selectedScanner !== ALL_SCANNERS) return
    try {
      if (selectedType === ITEM_TYPE_MISCONFIGURATION) {
        const response = await fetch(`/api/misconfigurations?resolved=${showResolved}&reportId=${selectedReport}${query ? "&search=" + query : ""}`)
        if (!response.ok) throw new Error(`Failed to fetch misconfigurations`);
        const data: Record<string, Misconfiguration[]> = await response.json();
        if (data === null || Object.keys(data).length === 0) {
          setMisconfigurations(null)
        } else {
          setMisconfigurations(data)
        }
      } else {
        const response = await fetch(`/api/vulnerabilities?reportId=${selectedReport}${query ? "&search=" + query : ""}`)
        if (!response.ok) throw new Error(`Failed to fetch vulnerabilities`);
        const data: Record<string, Vulnerability[]> = await response.json();
        if (data === null || Object.keys(data).length === 0) {
          setVulnerabilities(null)
        } else {
          setVulnerabilities(data)
        }
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let ws: WebSocket;
  
    fetch('/api/env')
      .then((res) => res.json())
      .then((data) => {
        ws = new WebSocket(data.websocketUrl);
  
        ws.onmessage = (event) => {
          const msg = event.data;
          if (msg.startsWith("@job")) {
            if (msg.endsWith("start")) setJobInProgress(true);
          } else if (msg.startsWith("@parse")) {
            if (msg.endsWith("end")) setJobInProgress(false);
          } else {
            setMessage(event.data);
          }
          setShowToast(true);
        };
      })
      .catch((err) => console.error("Failed to fetch WebSocket URL", err));
  
    return () => {
      if (ws) ws.close();
    };
  }, []);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const response = await fetch(`/api/reports?scanner=${selectedScanner}`);
        if (!response.ok) throw new Error(`Failed to fetch reports`);
        const data: Report[] = await response.json();
        setReports(data);
        if (data !== null && data.length > 0) {
          setSelectedReport(data[0].ID); 
        } else {
          setSelectedReport("")
          setMisconfigurations(null)
          setVulnerabilities(null)
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    if (selectedScanner === ALL_SCANNERS) {
      setSelectedReport("")
      setReports([])
      return
    }
    fetchReports();
  }, [selectedScanner, jobInProgress]);

  useEffect(() => {
    fetchData();
  }, [selectedReport, selectedType, query, showResolved]);

  useEffect(() => {
    setOpenCategory(null);
    setOpenItem(null);
  }, [selectedType])

  useEffect(() => {
    setOpenItem(null);
  }, [openCategory])

  const handleClick = () => {
    setShowConfirmationDialog(true)
    setConfirmationDialogAction(handleRescan)
  }

  const handleRescan = () => {
    setConfirmationDialogText(RERUN_SCAN_CONFIRMATION);
    setConfirmationDialogAction(() => async () => {
      try {
        const response = await fetch(`/api/run?scanner=${selectedScanner}`, {
          method: "POST",
        });
        if (!response.ok) throw new Error("Failed to trigger rescan");
        console.log("Rescan successful");
      } catch (err: any) {
        setMessage(err.message);
        setShowToast(true);
      }
      setShowConfirmationDialog(false);
    });
    setShowConfirmationDialog(true);
  };

  const handleDownload = () => {
    try {
      let jsonStr
      if (selectedType === ITEM_TYPE_MISCONFIGURATION) {
        jsonStr = JSON.stringify(misconfigurations, null, 2)
      } else {
        jsonStr = JSON.stringify(vulnerabilities, null, 2)
      }
      const blob = new Blob([jsonStr], { type: 'application/json' })

      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      const now = new Date()
      const timestamp = now.toLocaleString('sk-SK', {
        year: '2-digit',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      }).replaceAll('. ', '-').replaceAll(' ', '-').replaceAll(':', '-')
      a.download = `${selectedType.toLocaleLowerCase()}-export-${selectedScanner.replace(' ', '-').toLocaleLowerCase()}-${timestamp}.json`
      document.body.appendChild(a)
      a.click()
      a.remove()
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Download failed:', err)
    }
  }

  const handleDelete = () => {
    try {
      setConfirmationDialogText(DELETE_ALL_CONFIRMATION)
      setConfirmationDialogAction(() => async () => {
        if (selectedType === ITEM_TYPE_MISCONFIGURATION) {
          const response = await fetch(`/api/misconfigurations?reportId=${selectedReport}&search=${query}`, {
            method: "DELETE",
          });
          if (!response.ok) throw new Error("Failed to delete misconfigurations");
        } else {
          const response = await fetch(`/api/vulnerabilities?reportId=${selectedReport}&search=${query}`, {
            method: "DELETE",
          });
          if (!response.ok) throw new Error("Failed to delete vulnerabilities");
        }
        setShowConfirmationDialog(false)
        fetchData()
      })
      setShowConfirmationDialog(true)
    } catch (err: any) {
      setMessage(err.message);
      setShowToast(true);
    }
  };

  const handleResolve = () => {
    try {
      setConfirmationDialogText(RESOLVE_ALL_CONFIRMATION)
      setConfirmationDialogAction(() => async () => {
        if (selectedType === ITEM_TYPE_MISCONFIGURATION) {
          const response = await fetch(`/api/misconfigurations?reportId=${selectedReport}&search=${query}`, {
            method: "PATCH",
          });
          if (!response.ok) throw new Error("Failed to resolve misconfigurations");
          setMisconfigurations(prev => {
            if (prev == null || Object.keys(prev).length == 0) return prev;

            const newData = Object.fromEntries(
              Object.entries(prev).map(([category, misconfigs]) => [
                category,
                misconfigs.map((m) => ({
                  ...m,
                  Status:
                    m.Status === "RESOLVED"
                      ? "FAIL"
                      : m.Status === "FAIL" || m.Status === "MANUAL"
                      ? "RESOLVED"
                      : m.Status,
                })),
              ])
            );

            return newData
          })
        }
        setShowConfirmationDialog(false)
      })
      setShowConfirmationDialog(true)
    } catch (err: any) {
      setMessage(err.message);
      setShowToast(true);
    }
  };

  const handleResolveSingle = async (category: string, index: number, item: Misconfiguration) => {
    try {
      if (selectedType === ITEM_TYPE_MISCONFIGURATION && item.Status !== ITEM_STATUS_PASS) {
        const response = await fetch(`/api/misconfigurations?id=${item.ID}`, {
          method: "PATCH",
        });
        if (!response.ok) throw new Error("Failed to resolve misconfiguration");
        setMisconfigurations(prev => {
          if (prev == null) return null;

          const updatedCategory = [...prev[category]];
          const updatedItem = { ...updatedCategory[index], Status: updatedCategory[index].Status };
          updatedCategory[index] = updatedItem;

          return {
              ...prev,
              [category]: updatedCategory
          };
        });
        if (item.Status == ITEM_STATUS_FAIL || item.Status == ITEM_STATUS_MANUAL) {
          item.Status = ITEM_STATUS_RESOLVED;
        } else {
          item.Status = ITEM_STATUS_FAIL;
        }
      }
    } catch (err: any) {
      setMessage(err.message);
      setShowToast(true);
    }
  }

  const handleDeleteSingle = async (category: string, index: number, item: Item) => {
    try {
      if (selectedType === ITEM_TYPE_MISCONFIGURATION) {
        const response = await fetch(`/api/misconfigurations?id=${item.ID}`, {
          method: "DELETE",
        });
        if (!response.ok) throw new Error("Failed to delete misconfiguration");
        setMisconfigurations(prev => {
            if (prev == null) return null;
            
            const updatedCategory = [...prev[category]];
            updatedCategory.splice(index, 1);
            return {
                ...prev,
                [category]: updatedCategory
            };
        });
      } else {
        setVulnerabilities(prev => {
          if (prev == null) return null;
          
          const updatedCategory = [...prev[category]];
          updatedCategory.splice(index, 1);
          return {
              ...prev,
              [category]: updatedCategory
          };
      });
        const response = await fetch(`/api/vulnerabilities?id=${item.ID}`, {
          method: "DELETE",
        });
        if (!response.ok) throw new Error("Failed to delete vulnerability");
      }
    } catch (err: any) {
      setMessage(err.message);
      setShowToast(true);
    }
  }

  return (
    <div className="bg-slate-200 items-center justify-items-center min-h-screen p-8 gap-16">
      <h2 className="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">KSA Dashboard</h2>
      <FilterPanel 
        valueSelectedScanner={selectedScanner} onChangeSelectedScanner={setSelectedScanner} 
        valueSelectedType={selectedType} onChangeSelectedType={setSelectedType} 
        valueSelectedReport={selectedReport} onChangeSelectedReport={setSelectedReport} 
        valueReports={reports} 
        onClickRescan={handleClick}
        jobInProgress={jobInProgress}
        showResolved={showResolved}
        setShowResolved={setShowResolved}
        setLoading={setLoading}
        onChangeSearchValue={setQuery}
        onClickDownload={handleDownload}
        onClickDelete={handleDelete}
        onClickResolve={handleResolve}
      />
      {loading ? <div className="px-4"><Skeleton height={40} borderRadius={6} count={4} className="mt-4" baseColor="#1447e6"/></div> : error ? <div className="text-red-700">{error}</div> :
        selectedType === ITEM_TYPE_MISCONFIGURATION ? <MisconfigurationList data={misconfigurations} onClickDelete={handleDeleteSingle} onClickResolve={handleResolveSingle} openCategory={openCategory} setOpenCategory={setOpenCategory} openItem={openItem} setOpenItem={setOpenItem}/> : <VulnerabilityList data={vulnerabilities} onClickDelete={handleDeleteSingle} onClickResolve={undefined} openCategory={openCategory} setOpenCategory={setOpenCategory} openItem={openItem} setOpenItem={setOpenItem}/>}
      <SidePanel/>

      {showToast && (
              <Toast message={message} onClose={() => setShowToast(false)} />
            )}

      {showConfirmationDialog && (
          <ConfirmationDialog
            open={showConfirmationDialog}
            title="Are you sure?"
            description={confirmationDialogText}
            onConfirm={confirmationDialogAction}
            onCancel={() => setShowConfirmationDialog(false)}/>
        )}
    </div>
  );
}
