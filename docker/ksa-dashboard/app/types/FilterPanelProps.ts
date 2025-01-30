import { Report } from "./Report";

export type FilterPanelProps = {
    valueSelectedScanner: string;
    onChangeSelectedScanner: (newValue: string) => void;
    valueSelectedType: string;
    onChangeSelectedType: (newValue: string) => void;
    valueSelectedReport: string;
    onChangeSelectedReport: (newValue: string) => void;
    valueReports: Report[];
};