import os
import json
import datetime
from source import configuration


class DryRunHandler:
    """
    Handles dry-run functionality for the newsletter
    """
    
    def __init__(self):
        self.config = configuration.conf.dry_run
        self._ensure_output_directory()
    
    def _resolve_output_directory(self):
        """
        Resolve output directory based on environment and platform
        Returns the actual path where files should be saved
        """
        output_dir = self.config.output_directory
        
        # Handle Docker vs local environment
        if os.path.exists('/app'):
            # Docker environment - use absolute paths as configured
            if output_dir.startswith('/app'):
                return output_dir
            else:
                # If relative path, make it relative to /app
                return os.path.join('/app', output_dir.lstrip('./'))
        else:
            # Local environment - handle relative paths
            if os.path.isabs(output_dir):
                return output_dir
            else:
                # Make relative to current working directory
                return os.path.abspath(output_dir)
    
    def _ensure_output_directory(self):
        """Create output directory if it doesn't exist"""
        if self.config.enabled:
            actual_path = self._resolve_output_directory()
            try:
                os.makedirs(actual_path, exist_ok=True)
                # Log the actual path being used
                configuration.logging.info(f"Dry-run directory: {os.path.abspath(actual_path)}")
            except Exception as e:
                configuration.logging.error(f"Failed to create dry-run directory '{actual_path}': {e}")
                raise
    
    def _generate_filename(self, suffix=""):
        """Generate filename with date/timestamp placeholders"""
        now = datetime.datetime.now()
        
        filename = self.config.output_filename
        
        # Replace placeholders
        filename = filename.replace('{date}', now.strftime('%Y-%m-%d'))
        filename = filename.replace('{timestamp}', now.strftime('%Y%m%d_%H%M%S'))
        filename = filename.replace('{time}', now.strftime('%H%M%S'))
        
        # Add suffix if provided
        if suffix:
            name, ext = os.path.splitext(filename)
            filename = f"{name}_{suffix}{ext}"
        
        # Ensure we have the right directory
        output_dir = self._resolve_output_directory()
        return os.path.join(output_dir, filename)
    
    def _add_metadata_to_html(self, html_content, metadata):
        """Add metadata to HTML as comments"""
        if not self.config.include_metadata:
            return html_content
        
        metadata_comment = f"""<!--
Newsletter Generation Metadata:
Generated at: {metadata['generation_timestamp']}
Mode: {metadata['mode']}
SMTP Tested: {metadata['smtp_tested']}
Movies: {metadata['stats']['movies_count']}
TV Episodes: {metadata['stats']['tv_episodes_count']}
Email Size: {metadata['stats']['total_email_size_kb']}KB
-->"""
        
        return metadata_comment + '\n' + html_content
    
    def save_dry_run_output(self, html_content, metadata, mode="dry-run"):
        """
        Save HTML output and optional JSON metadata
        
        Args:
            html_content (str): The generated HTML email content
            metadata (dict): Email generation metadata
            mode (str): "dry-run" or "dry-run-smtp-only"
            
        Returns:
            tuple: (html_file_path, json_file_path or None)
        """
        if not self.config.enabled:
            return None, None
        
        try:
            # Generate filenames
            html_file = self._generate_filename()
            
            # Add metadata to HTML if enabled
            final_html = self._add_metadata_to_html(html_content, metadata)
            
            # Save HTML file
            with open(html_file, 'w', encoding='utf-8') as f:
                f.write(final_html)
            
            # Save JSON metadata if enabled
            json_file = None
            if self.config.save_email_data:
                json_file = self._generate_filename("data").replace('.html', '.json')
                with open(json_file, 'w', encoding='utf-8') as f:
                    json.dump(metadata, f, indent=2, ensure_ascii=False, default=str)
            
            # Log actual file locations
            configuration.logging.info(f"Dry-run output saved: {os.path.abspath(html_file)}")
            if json_file:
                configuration.logging.info(f"Metadata saved: {os.path.abspath(json_file)}")
            
            return html_file, json_file
            
        except Exception as e:
            configuration.logging.error(f"Failed to save dry-run files: {e}")
            raise
    
    def get_metadata(self, movies, series, total_tv, total_movie, mode="dry-run", smtp_tested=False):
        """Generate metadata for the email"""
        now = datetime.datetime.now()
        
        # Prepare movies data for JSON
        movies_list = []
        for movie_id, movie_data in movies.items():
            movies_list.append({
                "name": movie_data.get('name', 'Unknown'),
                "added_date": movie_data.get('created_on', '').split('T')[0],
                "tmdb_id": movie_data.get('tmdb_id', '')
            })
        
        # Prepare series data for JSON
        series_list = []
        for serie_id, serie_data in series.items():
            series_list.append({
                "series_name": serie_data.get('series_name', 'Unknown'),
                "seasons": serie_data.get('seasons', []),
                "episodes": serie_data.get('episodes', []),
                "added_date": serie_data.get('created_on', '').split('T')[0]
            })
        
        metadata = {
            "generation_timestamp": now.isoformat(),
            "mode": mode,
            "smtp_tested": smtp_tested,
            "jellyfin_server": configuration.conf.email_template.jellyfin_url,
            "stats": {
                "movies_count": total_movie,
                "tv_episodes_count": total_tv,
                "total_email_size_kb": 0  # Will be calculated after HTML generation
            },
            "movies": movies_list,
            "tv_shows": series_list,
            "recipients": configuration.conf.recipients if mode == "dry-run-smtp-only" else ["dry-run-mode"],
            "template_language": configuration.conf.email_template.language,
            "configuration_hash": str(hash(str(configuration.conf.__dict__)))
        }
        
        return metadata
