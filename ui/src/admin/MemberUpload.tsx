import { useCallback, useRef, useState } from 'react';
import { CloudUploadOutlined } from '@mui/icons-material';
import { Box, Button, Chip, Container, Paper, Stack, Typography } from '@mui/material';
import {
  DataGrid,
  GridColDef,
  GridToolbarContainer,
  GridToolbarQuickFilter,
} from '@mui/x-data-grid';
import * as XLSX from 'xlsx';
import { useSnackbar } from '../contexts/snackbar';

type MemberRow = {
  id: number;
  name: string;
  externalName: string;
  generation: string;
  email: string;
  errors: string[];
};

const EMAIL_TYPO_DOMAINS: Record<string, string[]> = {
  'gmail.com': ['gmail.con', 'gmail.cim', 'gmail.co', 'gmail.cm', 'gamil.com', 'gmial.com'],
  'naver.com': ['naver.con', 'naver.cim', 'naver.co', 'naver.cm', 'nave.com', 'navr.com'],
};

const validateName = (name: string): string | null => {
  if (!name || name.trim() === '') return 'Name is empty';
  if (/\d/.test(name)) return 'Name contains numbers';
  return null;
};

const validateGeneration = (gen: string): string | null => {
  if (!gen || gen.trim() === '') return 'Generation is empty';
  const num = Number(gen);
  if (isNaN(num)) return 'Generation is not a number';
  if (num <= 0) return 'Generation must be positive';
  return null;
};

const validateEmail = (email: string): string | null => {
  if (!email || email.trim() === '') return 'Email is empty';
  const basicPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!basicPattern.test(email)) return 'Invalid email format';

  const domain = email.split('@')[1]?.toLowerCase();
  for (const [correct, typos] of Object.entries(EMAIL_TYPO_DOMAINS)) {
    if (typos.includes(domain)) {
      return `Possible typo: did you mean ${correct}?`;
    }
  }
  return null;
};

const generateExternalNames = (names: string[]): string[] => {
  const counts = new Map<string, number>();
  for (const name of names) {
    counts.set(name, (counts.get(name) || 0) + 1);
  }

  const indices = new Map<string, number>();
  return names.map((name) => {
    if ((counts.get(name) || 0) > 1) {
      const idx = (indices.get(name) || 0) + 1;
      indices.set(name, idx);
      return `${name}${idx}`;
    }
    return name;
  });
};

const MemberUpload = () => {
  const { showWarning, showSuccess } = useSnackbar();
  const [rows, setRows] = useState<MemberRow[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const processFile = useCallback(
    (file: File) => {
      if (!file.name.endsWith('.csv') && !file.name.endsWith('.xlsx')) {
        showWarning('Please upload a CSV or XLSX file.');
        return;
      }

      const reader = new FileReader();
      reader.onload = (e) => {
        const data = e.target?.result;
        if (!data) return;

        const workbook = XLSX.read(data, { type: 'array' });
        const sheet = workbook.Sheets[workbook.SheetNames[0]];
        const json = XLSX.utils.sheet_to_json<Record<string, string>>(sheet, { raw: false });

        if (json.length === 0) {
          showWarning('The file is empty or has no data rows.');
          return;
        }

        const names = json.map((row) => (row.name || '').trim());
        const externalNames = generateExternalNames(names);

        const seen = new Set<string>();
        const duplicateKeys = new Set<string>();
        for (const row of json) {
          const key = `${(row.name || '').trim()}|${(row.generation || '').trim()}|${(row.email || '').trim()}`;
          if (seen.has(key)) duplicateKeys.add(key);
          seen.add(key);
        }

        const parsed: MemberRow[] = json.map((row, idx) => {
          const name = (row.name || '').trim();
          const generation = (row.generation || '').trim();
          const email = (row.email || '').trim();

          const errors: string[] = [];
          const rowKey = `${name}|${generation}|${email}`;
          if (duplicateKeys.has(rowKey)) errors.push('Duplicate row');
          const nameErr = validateName(name);
          if (nameErr) errors.push(`Name: ${nameErr}`);
          const genErr = validateGeneration(generation);
          if (genErr) errors.push(`Generation: ${genErr}`);
          const emailErr = validateEmail(email);
          if (emailErr) errors.push(`Email: ${emailErr}`);

          return {
            id: idx + 1,
            name,
            externalName: externalNames[idx],
            generation,
            email,
            errors,
          };
        });

        setRows(parsed);
        const errorCount = parsed.filter((r) => r.errors.length > 0).length;
        if (errorCount === 0) {
          showSuccess(`All ${parsed.length} rows are valid.`);
        }
      };
      reader.readAsArrayBuffer(file);
    },
    [showWarning, showSuccess],
  );

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) processFile(file);
    e.target.value = '';
  };

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragOver(false);
      const file = e.dataTransfer.files[0];
      if (file) processFile(file);
    },
    [processFile],
  );

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = () => {
    setIsDragOver(false);
  };

  const handleDownloadTemplate = () => {
    const ws = XLSX.utils.aoa_to_sheet([['name', 'generation', 'email']]);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Template');
    XLSX.writeFile(wb, 'member_upload_template.csv', { bookType: 'csv' });
  };

  const handleClear = () => {
    setRows([]);
  };

  const errorCount = rows.filter((r) => r.errors.length > 0).length;
  const validCount = rows.length - errorCount;

  const columns: GridColDef[] = [
    {
      field: 'name',
      headerName: 'Name',
      width: 140,
      renderCell: (params) => {
        const hasNameError = params.row.errors.some((e: string) => e.startsWith('Name:'));
        return <span style={hasNameError ? { color: '#d32f2f' } : undefined}>{params.value}</span>;
      },
    },
    {
      field: 'externalName',
      headerName: 'External Name',
      width: 150,
    },
    {
      field: 'generation',
      headerName: 'Generation',
      width: 120,
      renderCell: (params) => {
        const hasGenError = params.row.errors.some((e: string) => e.startsWith('Generation:'));
        return <span style={hasGenError ? { color: '#d32f2f' } : undefined}>{params.value}</span>;
      },
    },
    {
      field: 'email',
      headerName: 'Email',
      width: 240,
      renderCell: (params) => {
        const hasEmailError = params.row.errors.some((e: string) => e.startsWith('Email:'));
        return <span style={hasEmailError ? { color: '#d32f2f' } : undefined}>{params.value}</span>;
      },
    },
    {
      field: 'status',
      headerName: 'Status',
      flex: 1,
      minWidth: 200,
      renderCell: (params) => {
        const errors: string[] = params.row.errors;
        if (errors.length === 0) {
          return <Chip label="Valid" color="success" size="small" />;
        }
        return (
          <Box sx={{ py: 1 }}>
            {errors.map((err) => (
              <Typography key={err} variant="body2" color="error" sx={{ lineHeight: 1.4 }}>
                {err}
              </Typography>
            ))}
          </Box>
        );
      },
      sortComparator: (_v1, _v2, param1, param2) => {
        const e1 = (param1.api.getRow(param1.id) as MemberRow).errors.length;
        const e2 = (param2.api.getRow(param2.id) as MemberRow).errors.length;
        return e2 - e1;
      },
    },
  ];

  return (
    <Container>
      <Paper sx={{ p: 2, mt: 2 }} elevation={4}>
        <Stack direction="row" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6">Member Upload</Typography>
          <Button variant="outlined" size="small" onClick={handleDownloadTemplate}>
            Download Template
          </Button>
        </Stack>

        <Paper
          variant="outlined"
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onClick={() => fileInputRef.current?.click()}
          sx={{
            p: 4,
            mb: 2,
            textAlign: 'center',
            cursor: 'pointer',
            border: isDragOver ? '2px dashed' : '1px dashed',
            borderColor: isDragOver ? 'primary.main' : 'divider',
            bgcolor: isDragOver ? 'action.hover' : 'transparent',
            transition: 'all 0.2s',
          }}
        >
          <CloudUploadOutlined sx={{ fontSize: 40, color: 'text.secondary', mb: 1 }} />
          <Typography variant="body1" color="text.secondary">
            Drag & drop a CSV file here, or click to browse
          </Typography>
          <Typography variant="caption" color="text.disabled">
            Accepts .csv and .xlsx files
          </Typography>
          <input
            ref={fileInputRef}
            type="file"
            accept=".csv,.xlsx"
            hidden
            onChange={handleFileChange}
          />
        </Paper>

        {rows.length > 0 && (
          <>
            <Stack direction="row" spacing={1} alignItems="center" mb={2}>
              <Chip label={`Total: ${rows.length}`} variant="outlined" size="small" />
              <Chip label={`Valid: ${validCount}`} color="success" size="small" />
              {errorCount > 0 && (
                <Chip label={`Errors: ${errorCount}`} color="error" size="small" />
              )}
              <Box sx={{ flexGrow: 1 }} />
              <Button variant="text" size="small" onClick={handleClear}>
                Clear
              </Button>
            </Stack>

            <Box sx={{ height: 'calc(100vh - 380px)', width: '100%' }}>
              <DataGrid
                rows={rows}
                columns={columns}
                getRowHeight={() => 'auto'}
                disableRowSelectionOnClick
                slots={{
                  toolbar: () => (
                    <GridToolbarContainer sx={{ pt: 2, px: 2 }}>
                      <GridToolbarQuickFilter />
                    </GridToolbarContainer>
                  ),
                }}
                initialState={{
                  sorting: {
                    sortModel: [{ field: 'status', sort: 'desc' }],
                  },
                  pagination: {
                    paginationModel: { pageSize: 20, page: 0 },
                  },
                }}
                sx={{
                  '& .MuiDataGrid-cell': {
                    alignItems: 'center',
                  },
                }}
              />
            </Box>
          </>
        )}
      </Paper>
    </Container>
  );
};

export default MemberUpload;
