import { create } from "zustand";

type DashboardState = {
  selectedTeamId: string | null;
  setSelectedTeamId: (teamId: string | null) => void;
};

export const useDashboardStore = create<DashboardState>((set) => ({
  selectedTeamId: null,
  setSelectedTeamId: (teamId) => set({ selectedTeamId: teamId }),
}));
