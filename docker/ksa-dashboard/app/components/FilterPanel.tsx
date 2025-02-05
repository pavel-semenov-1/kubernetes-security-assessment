'use client';
import React from "react";
import { Report } from "../types/Report";
import { RotateCcw } from "lucide-react";
import SearchInput from "./SearchInput";

interface FilterPanelProps {
    valueSelectedScanner: string;
    onChangeSelectedScanner: (newValue: string) => void;
    valueSelectedType: string;
    onChangeSelectedType: (newValue: string) => void;
    valueSelectedReport: string;
    onChangeSelectedReport: (newValue: string) => void;
    valueReports: Report[];
    onClickRescan: () => void;
    jobInProgress: boolean;
    setLoading: (newValue: boolean) => void;
    onChangeSearchValue: (newValue: string) => void;
};

const FilterPanel = ({ valueSelectedScanner, onChangeSelectedScanner, valueSelectedType, onChangeSelectedType, valueSelectedReport, onChangeSelectedReport, valueReports, onClickRescan, jobInProgress, setLoading, onChangeSearchValue }: FilterPanelProps) => {
    const scanners = ["Trivy", "Kubescape"];
    const types = ["Misconfiguration", "Vulnerability"];

    return (
        <div className="flex items-center gap-4 p-4">
            <div className="flex items-center gap-3">
                <label htmlFor="scanner" className="text-base text-gray-800">
                    Scanner:
                </label>
                <select
                    id="scanner"
                    value={valueSelectedScanner}
                    onChange={(e) => onChangeSelectedScanner(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                >
                    {scanners.map((scanner, index) => (
                        <option key={index} value={scanner}>
                            {scanner}
                        </option>
                    ))}
                </select>
            </div>

            <div className="flex items-center gap-3">
                <label htmlFor="type" className="text-base text-gray-800">
                    Type:
                </label>
                <select
                    id="type"
                    value={valueSelectedType}
                    onChange={(e) => onChangeSelectedType(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {types.map((type, index) => (
                        <option key={index} value={type}>
                        {type}
                        </option>
                    ))}
                    </select>
            </div>

            <div className="flex items-center gap-3">
                <label htmlFor="report" className="text-base text-gray-800">
                    Report:
                </label>
                {valueReports.length > 0 ?
                <select
                    id="report"
                    value={valueSelectedReport}
                    onChange={(e) => onChangeSelectedReport(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {valueReports.map((report) => (
                        <option key={report.ID} value={report.ID}>
                        {report.Filename}
                        </option>
                    ))}
                </select>
                : 
                <span>No data. Trigger a scan first.</span>
                }
            </div>

            <div className="gap-3 bg-blue-700 text-white border rounded-full p-1">
                <a onClick={onClickRescan} className="button"><RotateCcw className={jobInProgress ? "animate-spin" : ""}/></a>
            </div>

            <SearchInput setLoading={setLoading} setQuery={onChangeSearchValue} className="relative ml-auto"/>
        </div>
    );
}

export default FilterPanel;