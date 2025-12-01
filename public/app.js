(() => {
  const apiBase = '/api/licenseplate';

  const $ = id => document.getElementById(id);
  const tbody = document.querySelector('#records-table tbody');
  const info = $('records-info');

  async function loadRecords(query) {
    try {
      let url = apiBase + '/records';
      if (query) url += '?q=' + encodeURIComponent(query);
      const res = await fetch(url);
      const data = await res.json();
      // API returns {count: N, records: [...]}
      renderRecords(data.records || []);
    } catch (err) {
      console.error(err);
      info.textContent = 'Failed to load records';
    }
  }

  function renderRecords(records) {
    tbody.innerHTML = '';
    info.textContent = `${records.length} record(s)`;
    records.forEach(r => {
      const isExpired = r.access_expires_at && new Date(r.access_expires_at) < new Date();
      const expiryClass = isExpired ? 'expired' : '';
      const typeBadge = getTypeBadge(r.visitor_type || 'guest');
      
      const tr = document.createElement('tr');
      if (isExpired) tr.classList.add('expired-row');
      
      tr.innerHTML = `
        <td>${escapeHtml(r.plate_number)}</td>
        <td>${escapeHtml(r.guest_name || '')}</td>
        <td>${typeBadge}</td>
        <td>${escapeHtml(r.room_number || '')}</td>
        <td>${formatDate(r.check_in)}</td>
        <td>${formatDate(r.check_out)}</td>
        <td class="${expiryClass}">${formatDate(r.access_expires_at) || '∞'}</td>
        <td>${escapeHtml(r.notes || '')}</td>
        <td>
          <button class="btn" data-plate="${escapeHtml(r.plate_number)}" data-action="view">View</button>
          <button class="btn danger" data-plate="${escapeHtml(r.plate_number)}" data-action="delete">Delete</button>
        </td>
      `;
      tbody.appendChild(tr);
    });
  }

  function getTypeBadge(type) {
    const badges = {
      guest: '<span class="badge badge-guest">Guest</span>',
      visitor: '<span class="badge badge-visitor">Visitor</span>',
      staff: '<span class="badge badge-staff">Staff</span>',
      delivery: '<span class="badge badge-delivery">Delivery</span>',
      contractor: '<span class="badge badge-contractor">Contractor</span>',
      vip: '<span class="badge badge-vip">VIP</span>',
    };
    return badges[type] || badges.guest;
  }

  function escapeHtml(s){
    if(!s) return '';
    return s.replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"})[c]);
  }

  function formatDate(v){
    if(!v) return '';
    try{ return new Date(v).toLocaleString(); }catch(e){return v}
  }

  // Scan form
  $('scan-form').addEventListener('submit', async (ev) => {
    ev.preventDefault();
    const expiryInput = $('access_expires_at').value;
    const payload = {
      plate_number: $('plate_number').value.trim(),
      guest_name: $('guest_name').value.trim(),
      room_number: $('room_number').value.trim(),
      vehicle_make: $('vehicle_make').value.trim(),
      vehicle_model: $('vehicle_model').value.trim(),
      visitor_type: $('visitor_type').value,
      purpose: $('purpose').value.trim(),
    };
    
    // Convert datetime-local to ISO 8601 if provided
    if (expiryInput) {
      payload.access_expires_at = new Date(expiryInput).toISOString();
    }
    
    try{
      const res = await fetch(apiBase + '/scan', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      if(!res.ok) throw new Error('Failed to save');
      $('scan-form').reset();
      await loadRecords();
      info.textContent = 'Saved successfully';
    }catch(err){
      console.error(err);
      info.textContent = 'Save failed';
    }
  });

  // Table actions
  tbody.addEventListener('click', async (ev) => {
    const btn = ev.target.closest('button');
    if(!btn) return;
    const plate = btn.dataset.plate;
    const action = btn.dataset.action;
    if(action === 'delete'){
      if(!confirm('Delete record for ' + plate + '?')) return;
      try{
        const res = await fetch(apiBase + '/records/' + encodeURIComponent(plate), {method:'DELETE'});
        if(!res.ok) throw new Error('Delete failed');
        await loadRecords();
      }catch(err){
        console.error(err);
        info.textContent = 'Delete failed';
      }
    }
    if(action === 'view'){
      try{
        const res = await fetch(apiBase + '/records/' + encodeURIComponent(plate));
        if(!res.ok) throw new Error('Not found');
        const r = await res.json();
        alert(JSON.stringify(r, null, 2));
      }catch(err){
        console.error(err);
        info.textContent = 'Failed to load record';
      }
    }
  });

  // Search / refresh
  $('btn-refresh').addEventListener('click', () => loadRecords($('search').value.trim()));
  $('search').addEventListener('keyup', (e) => { if(e.key === 'Enter') loadRecords($('search').value.trim()); });

  // initial load
  loadRecords();

})();
