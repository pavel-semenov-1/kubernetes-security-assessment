import { SeverityLevels } from "./SeverityLevels";

export type Stats = {
  vulnerabilities: Record<SeverityLevels, number>;
  misconfigurations: Record<SeverityLevels, number>;
};