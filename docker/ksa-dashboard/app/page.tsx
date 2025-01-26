'use client';
import FilterPanel from "./components/FilterPanel";
import React, { useEffect, useState } from "react";
import { Misconfiguration } from "./types/Misconfiguration";
import { ItemListProps } from "./types/ItemListProps";
import { Vulnerability } from "./types/Vulnerability";
import MisconfigList from "./components/MisconfigList";
import VulnList from "./components/VulnList";

export default function Home() {
  const categories = ["CRITICAL", "HIGH", "MEDIUM", "LOW"];
  const [type, setType] = useState("Misconfiguration")
  const [scanner, setScanner] = useState("Trivy")
  const [misconfigurations, setMisconfigurations] = useState<ItemListProps<Misconfiguration>>({
    data: categories.reduce((acc, category) => {
      acc[category] = [];
      return acc;
    }, {} as Record<string, Misconfiguration[]>)
  });
  const [vulnerabilities, setVulnerabilities] = useState<ItemListProps<Vulnerability>>({
    data: categories.reduce((acc, category) => {
      acc[category] = [];
      return acc;
    }, {} as Record<string, Vulnerability[]>)
  });
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMisconfigurations = async () => {
      try {
        const results = await Promise.all(
          categories.map(async (category) => {
            const response = await fetch(`/api/misconfigurations?scanner=${scanner}&severity=${category}`);
            if (!response.ok) {
              throw new Error(`Failed to fetch misconfigurations for ${category}`);
            }
            const data: Misconfiguration[] = await response.json();
            return { category, data };
          })
        );

        const mappedData: Record<string, Misconfiguration[]> = results.reduce((acc, { category, data }) => {
          acc[category] = data;
          return acc;
        }, {} as Record<string, Misconfiguration[]>);


        setMisconfigurations({ data: mappedData });
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    const fetchVulnerabilities = async () => {
      try {
        const results = await Promise.all(
          categories.map(async (category) => {
            const response = await fetch(`/api/vulnerabilities?scanner=${scanner}&severity=${category}`);
            if (!response.ok) {
              throw new Error(`Failed to fetch vulnerabilities for ${category}`);
            }
            const data: Vulnerability[] = await response.json();
            return { category, data };
          })
        );

        const mappedData: Record<string, Vulnerability[]> = results.reduce((acc, { category, data }) => {
          acc[category] = data;
          return acc;
        }, {} as Record<string, Vulnerability[]>);


        setVulnerabilities({ data: mappedData });
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };


    fetchMisconfigurations();
    fetchVulnerabilities();
  }, [categories, scanner]);

  return (
    <div className="bg-slate-200 items-center justify-items-center min-h-screen p-8 gap-16">
      <h2 className="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">KSA Dashboard</h2>
      <FilterPanel valueScanner={scanner} onChangeScanner={setScanner} valueType={type} onChangeType={setType}/>
      { loading ? <div>Loading...</div> : 
        error ? <div className="text-red-700">{error}</div> :
      ( type === "Misconfiguration" ? <MisconfigList data={misconfigurations.data}/> : <VulnList data={vulnerabilities.data}/>)}
    </div>
  );
}
