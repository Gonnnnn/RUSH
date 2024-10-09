import { Box, Button, CircularProgress, Grid, Paper, Typography } from '@mui/material';
import { QRCodeCanvas } from 'qrcode.react';
import { Session } from '../client/http';
import { formatDateToMonthDate } from '../common/date';

const AttendanceQrPanel = ({
  session,
  qrRef,
  qrSizePx,
  qrDownloadSizePx,
  isCreatingForm,
  onCreateQRCode,
}: {
  session: Session;
  qrRef: React.RefObject<HTMLDivElement>;
  qrSizePx: number;
  qrDownloadSizePx: number;
  isCreatingForm: boolean;
  onCreateQRCode: () => void;
}) => (
  <Paper sx={{ p: 2 }} elevation={4}>
    {session.googleFormUri ? (
      <>
        <Typography variant="h6" gutterBottom>
          출석 QR
        </Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2, my: 2 }}>
          <QRCodeCanvas value={session.googleFormUri} size={qrSizePx} />
          {/* Make a hidden QR for download. The QR for display is too small that it breaks when resizing for downloading. */}
          <div ref={qrRef} style={{ display: 'None' }}>
            <QRCodeCanvas value={session.googleFormUri} size={qrDownloadSizePx} />
          </div>
          <Grid container spacing={2} justifyContent="center">
            <Grid item xs={12} sm={6}>
              <Button
                variant="outlined"
                fullWidth
                // to Month Day format
                onClick={() => onQrDownload(qrRef, qrDownloadSizePx, formatDateToMonthDate(new Date(session.startsAt)))}
              >
                QR 다운로드
              </Button>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Button variant="outlined" fullWidth onClick={() => window.open(session.googleFormUri, '_blank')}>
                Google form 열기 (제출용)
              </Button>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Button
                variant="outlined"
                fullWidth
                onClick={() => window.open(`https://docs.google.com/forms/d/${session.googleFormId}/edit`, '_blank')}
              >
                Google form 열기 (편집용)
              </Button>
            </Grid>
          </Grid>
        </Box>
      </>
    ) : (
      <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2, my: 2 }}>
        <Typography variant="h6">No google form attached yet!</Typography>
        <Button variant="contained" onClick={onCreateQRCode} disabled={isCreatingForm}>
          {isCreatingForm ? <CircularProgress size={24} /> : 'Create QR code'}
        </Button>
      </Box>
    )}
  </Paper>
);

const onQrDownload = (qrRef: React.RefObject<HTMLDivElement>, qrSize: number, text: string) => {
  if (!qrRef.current) return;

  const canvas = qrRef.current.querySelector('canvas');
  if (!canvas) return;

  const newCanvas = document.createElement('canvas');
  const ctx = newCanvas.getContext('2d');
  if (!ctx) return;

  const paddingPx = 64;
  const textSpacePx = 128;
  const newCanvasWidth = qrSize + paddingPx * 2;
  const newCanvasHeight = qrSize + paddingPx * 2 + textSpacePx;

  newCanvas.width = newCanvasWidth;
  newCanvas.height = newCanvasHeight;
  ctx.fillStyle = 'white';
  ctx.fillRect(0, 0, newCanvas.width, newCanvas.height);

  const qrYoffset = (newCanvasHeight - qrSize) / 2;
  const qrXoffset = (newCanvasWidth - qrSize) / 2;
  ctx.drawImage(canvas, qrXoffset, qrYoffset, qrSize, qrSize);

  ctx.font = '32px Helvetica';
  ctx.fillStyle = 'black';
  ctx.textAlign = 'center';
  ctx.fillText(text, newCanvasWidth / 2, Math.min(qrYoffset + qrSize + textSpacePx / 2, newCanvasHeight - 16));

  const a = document.createElement('a');
  a.href = newCanvas.toDataURL('image/png');
  a.download = `${text.replace(/ /g, '_')}.png`;
  a.click();
};

export default AttendanceQrPanel;
