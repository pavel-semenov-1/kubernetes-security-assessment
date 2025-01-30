'use client';
import React from "react";
import { FilterPanelProps } from "../types/FilterPanelProps";

const FilterPanel = ({ valueSelectedScanner, onChangeSelectedScanner, valueSelectedType, onChangeSelectedType, valueSelectedReport, onChangeSelectedReport, valueReports }: FilterPanelProps) => {
    const scanners = ["Trivy", "Kubescape"];
    const types = ["Misconfiguration", "Vulnerability"];

    return (
        <div className="flex items-center gap-4 p-4 max-w-2xl">
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
                <select
                    id="report"
                    value={valueSelectedReport}
                    onChange={(e) => onChangeSelectedReport(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {valueReports.map((report) => (
                        <option key={report.Id} value={report.Filename}>
                        {report.Filename}
                        </option>
                    ))}
                    </select>
            </div>
        </div>
    );
}

export default FilterPanel;