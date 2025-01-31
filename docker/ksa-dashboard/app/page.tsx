'use client';
import FilterPanel from "./components/FilterPanel";
import React, { useEffect, useState } from "react";
import { Misconfiguration } from "./types/Misconfiguration";
import { ItemListProps } from "./types/ItemListProps";
import { Vulnerability } from "./types/Vulnerability";
import MisconfigurationList from "./components/MisconfigurationList";
import VulnerabilityList from "./components/VulnerabilityList";
import { Report } from "./types/Report";
import Toast from "./components/Toast";
import SidePanel from "./components/SidePanel";
import ConfirmationDialog from "./components/ConfirmationDialog";

const categories = ["CRITICAL", "HIGH", "MEDIUM", "LOW"];

export default function Home() {
  const [selectedType, setSelectedType] = useState("Misconfiguration");
  const [selectedScanner, setSelectedScanner] = useState("Trivy");
  const [selectedReport, setSelectedReport] = useState("");
  const [reports, setReports] = useState<Report[]>([]);
  const [misconfigurations, setMisconfigurations] = useState<ItemListProps<Misconfiguration>>({
    data: { CRITICAL: [], HIGH: [], MEDIUM: [], LOW: [] },
  });
  const [vulnerabilities, setVulnerabilities] = useState<ItemListProps<Vulnerability>>({
    data: { CRITICAL: [], HIGH: [], MEDIUM: [], LOW: [] },
  });
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [showToast, setShowToast] = useState(false);
  const [message, setMessage] = useState("");
  const [jobInProgress, setJobInProgress] = useState(false);
  const [showConfirmationDialog, setShowConfirmationDialog] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(`${process.env.NEXT_PUBLIC_AGGREGATOR_API_URL?.replace("http", "ws")}/ws`);

    ws.onmessage = (event) => {
      const msg = event.data;
      if (msg.startsWith("@")) {
        if (msg.endsWith("start")) {
          setJobInProgress(true)
        } else if (msg.endsWith("end")) {
          setJobInProgress(false)
        }
      } else {
        setMessage(event.data);
      }
      setShowToast(true);
    };

    return () => ws.close();
  }, []);

  useEffect(() => {
    const fetchReports = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/api/reports?scanner=${selectedScanner}`);
        if (!response.ok) throw new Error(`Failed to fetch reports`);
        const data: Report[] = await response.json();
        setReports(data);
        if (!selectedReport && data.length > 0) {
          setSelectedReport(data[0].ID); 
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchReports();
  }, [selectedScanner, jobInProgress]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [misconfigResults, vulnResults] = await Promise.all([
          Promise.all(categories.map(async (category) => {
            const response = await fetch(`/api/misconfigurations?scanner=${selectedScanner}&severity=${category}&reportId=${selectedReport}`);
            if (!response.ok) throw new Error(`Failed to fetch misconfigurations for ${category}`);
            const data: Misconfiguration[] = await response.json();
            return { category, data };
          })),
          Promise.all(categories.map(async (category) => {
            const response = await fetch(`/api/vulnerabilities?scanner=${selectedScanner}&severity=${category}&reportId=${selectedReport}`);
            if (!response.ok) throw new Error(`Failed to fetch vulnerabilities for ${category}`);
            const data: Vulnerability[] = await response.json();
            return { category, data };
          })),
        ]);

        setMisconfigurations({
          data: misconfigResults.reduce((acc, { category, data }) => ({ ...acc, [category]: data }), {}),
        });
        setVulnerabilities({
          data: vulnResults.reduce((acc, { category, data }) => ({ ...acc, [category]: data }), {}),
        });
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    if (!selectedReport) return;
    fetchData();
  }, [selectedReport]);

  const handleClick = async () => {
    setShowConfirmationDialog(true)
  }

  const handleRescan = async () => {
    try {
      setShowConfirmationDialog(false)
      const response = await fetch(`/api/run?scanner=${selectedScanner}`, {
        method: "POST",
      });
      if (!response.ok) throw new Error("Failed to trigger rescan");
      console.log("Rescan successful");
    } catch (err: any) {
      setMessage(err.message);
      setShowToast(true);
    }
  };

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
      />
      {loading ? <div className="text-black">Loading...</div> : error ? <div className="text-red-700">{error}</div> :
        selectedType === "Misconfiguration" ? <MisconfigurationList data={misconfigurations.data}/> : <VulnerabilityList data={vulnerabilities.data}/>}
      <SidePanel/>

      {showToast && (
              <Toast message={message} onClose={() => setShowToast(false)} />
            )}

      {showConfirmationDialog && (
          <ConfirmationDialog
            open={showConfirmationDialog}
            title="Are you sure?"
            description="This would run a Kubernetes job and generate a new report."
            onConfirm={handleRescan}
            onCancel={() => setShowConfirmationDialog(false)}/>
        )}
    </div>
  );
}
